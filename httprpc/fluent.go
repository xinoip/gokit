package httprpc

import "github.com/xinoip/gokit/httpapi"

// Handle registers an RPC handler as an HTTP POST endpoint.
func Handle[TReq, TRes any](
	r *httpapi.Registry,
	handler httpapi.Handler[TReq, TRes],
	name string,
	opts ...httpapi.Option,
) {
	rpcHandler := MakeHandler(handler)
	httpapi.Post(r, "/rpc/"+name, rpcHandler, name, opts...)
}
