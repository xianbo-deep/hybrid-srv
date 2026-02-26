package httpx

import (
	"encoding/json"
	"hybrid-srv/core"
	"net/http"
)

func httpStatusFromBizCode(code int) int {
	switch code {
	case core.CodeOK:
		return http.StatusOK
	case core.CodeBadRequest:
		return http.StatusBadRequest
	case core.CodeUnauthorized:
		return http.StatusUnauthorized
	case core.CodeForbidden:
		return http.StatusForbidden
	case core.CodeNotFound:
		return http.StatusNotFound
	case core.CodeInternal:
		return http.StatusInternalServerError
	default:
		// 兜底：业务失败但没分类 -> 500
		if code != 0 {
			return http.StatusInternalServerError
		}
		return http.StatusOK
	}
}

func renderResultHTTP(w *core.ResponseWriterWrapper, r core.Result) {
	for k, v := range r.Meta {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 写 header
	w.WriteHeader(httpStatusFromBizCode(r.Code))

	_ = json.NewEncoder(w).Encode(r)
}
