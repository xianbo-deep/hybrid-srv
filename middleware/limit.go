package middleware

import (
	"container/list"
	"net/http"
	"sync"

	"github.com/xianbo-deep/Fuse/core"
	"golang.org/x/time/rate"
)

type LRULimiter struct {
	capacity int                      // 容量
	mu       sync.Mutex               // 锁
	lrulist  *list.List               // 双向链表
	items    map[string]*list.Element // 哈希表 实现O(1)查找
}

func NewLRULimiter(capacity int) *LRULimiter {
	return &LRULimiter{
		capacity: capacity,
		lrulist:  list.New(),
		items:    make(map[string]*list.Element),
	}
}

type entry struct {
	ip      string
	limiter *rate.Limiter
}

func (l *LRULimiter) GetLimiter(clientIP string, tokens, burst int) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 命中缓存
	if elem, ok := l.items[clientIP]; ok {
		l.lrulist.MoveToFront(elem)
		return elem.Value.(*entry).limiter
	}

	// 未命中 新建
	limiter := rate.NewLimiter(rate.Limit(tokens), burst)
	newEntry := &entry{ip: clientIP, limiter: limiter}

	// 新节点插入头部
	elem := l.lrulist.PushFront(newEntry)
	l.items[clientIP] = elem

	// 查看是否超过最大容量
	if l.lrulist.Len() > l.capacity {
		// 删除最久未使用节点
		oldest := l.lrulist.Back()
		if oldest != nil {
			l.lrulist.Remove(oldest)
			delete(l.items, oldest.Value.(*entry).ip)
		}
	}

	return limiter
}

// RateLimiterConfig 限流器配置。
type RateLimiterConfig struct {
	Tokens int // 每秒生成几个令牌
	Burst  int // 桶的容量大小 决定可以同时处理几个请求
}

func NewRateLimiterConfig(tokens, burst int) *RateLimiterConfig {
	return &RateLimiterConfig{Tokens: tokens, Burst: burst}
}

var (
	lruLimiter               = NewLRULimiter(10000)
	defaultRateLimiterConfig = RateLimiterConfig{
		Tokens: 10,
		Burst:  1,
	}
)

// RateLimit 限流中间件。
//
// 使用 LRU 缓存缓存最近的有限个数的限流器，防止限流器数量太多导致 OOM 。
func RateLimit(config ...RateLimiterConfig) core.HandlerFunc {
	return func(c core.Ctx) core.Result {
		var cfg = defaultRateLimiterConfig
		if len(config) > 0 {
			cfg.Burst = config[0].Burst
			cfg.Tokens = config[0].Tokens
		}
		clientIP := c.ClientIP()

		limiter := lruLimiter.GetLimiter(clientIP, cfg.Tokens, cfg.Burst)
		if !limiter.Allow() {
			return core.Fail(core.CodeBadRequest, "Too Many Requests").WithHttpStatus(http.StatusTooManyRequests)
		}
		return c.Next()
	}
}
