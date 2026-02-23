package core

import (
	"context"
	"time"
)

type Ctx struct {
	// 继承
	context.Context

	// 通用字段
	RequestID string
	StartTime time.Time

	// 键值表
	Values map[string]any

	// 中止标志
	aborted bool
}

// 创建有基础信息的上下文
func NewCtx(parent context.Context) *Ctx {
	if parent == nil {
		parent = context.Background()
	}
	return &Ctx{
		Context:   parent,
		StartTime: time.Now(),
		Values:    make(map[string]any),
	}
}

func (ctx *Ctx) Abort() {
	ctx.aborted = true
}

func (ctx *Ctx) Aborted() bool {
	return ctx.aborted
}
