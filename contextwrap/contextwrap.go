package contextwrap

import (
	"context"
	"net/http"

	response "github.com/danielpnjt/go-library/basic"
)

const (
	traceKey      = "trace"
	bodyKey       = "body"
	thirdPartyKey = "thirdParty"
	respKey       = "resp"
)

func SetTraceFromContext(ctx context.Context, trace []interface{}) context.Context {
	ctx = context.WithValue(ctx, traceKey, trace)
	return ctx
}

func GetTraceFromContext(ctx context.Context) []interface{} {
	lr := ctx.Value(traceKey)
	if l, ok := lr.([]interface{}); ok {
		return l
	} else {
		return []interface{}{}
	}
}

func GetBody(r *http.Request) []byte {
	lr := r.Context().Value(bodyKey)
	if l, ok := lr.([]byte); ok {
		return l
	} else {
		return []byte("")
	}
}

func GetBodyFromContext(ctx context.Context) []byte {
	lr := ctx.Value(bodyKey)
	if l, ok := lr.([]byte); ok {
		return l
	} else {
		return []byte("")
	}
}

func GetThirdPartyFromContext(ctx context.Context) string {
	lr := ctx.Value(thirdPartyKey)
	if l, ok := lr.(string); ok {
		return l
	} else {
		return ""
	}
}

func SetThirdPartyFromContext(ctx context.Context, thirdParty string) context.Context {
	ctx = context.WithValue(ctx, thirdPartyKey, thirdParty)
	return ctx
}

func GetResponseFromContext(ctx context.Context) *response.Response {
	lr := ctx.Value(respKey)
	if lr == nil {
		return &response.Response{}
	}
	if l, ok := lr.(*response.Response); ok {
		return l
	} else {
		return &response.Response{}
	}
}
