package httprpc

import "github.com/xinoip/gokit/api"

// RPC is same as [api.POST] but intentionally limited to create RPC endpoints.
func RPC[TReq, TRes any](r *api.Registry, handler api.Handler[TReq, TRes], name string, opts ...api.Option) {
	rpcHandler := MakeHandler(handler)
	api.Post(r, "/rpc/"+name, rpcHandler, name, opts...)
}
