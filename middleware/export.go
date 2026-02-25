package middleware

import (
	"hybrid-srv/core"
)

func Defaults() []core.HandlerFunc {
	return []core.HandlerFunc{
		Recovery(),
		RequestID(),
		Logger(),
	}
}
