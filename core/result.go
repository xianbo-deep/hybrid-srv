package core

/*
跨协议返回值
*/
type Result struct {
	// 0表示成功 其它表示失败
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`

	// 跨协议元信息
	Meta map[string]string `json:"meta,omitempty"`
}

func OK(data any) Result {
	return Result{
		Code: 0,
		Data: data,
	}
}

func OKmsg(msg string, data any) Result {
	return Result{
		Code: 0,
		Data: data,
		Msg:  msg,
	}
}

func Fail(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
	}
}

func (r Result) WithMeta(k, v string) Result {
	if r.Meta == nil {
		r.Meta = map[string]string{}
	}
	r.Meta[k] = v
	return r
}
