package api

// Option is a modification for an endpoint. Can be used to further customize an endpoint.
type Option func(*Endpoint)

// WithInsecure marks the endpoint as insecure.
func WithInsecure() Option {
	return func(e *Endpoint) {
		e.Secure = false
	}
}

// WithMiddlewares adds middlewares specific to the Endpoint.
func WithMiddlewares(m ...Middleware) Option {
	return func(e *Endpoint) {
		for _, middleware := range m {
			e.Middlewares = append(e.Middlewares, middleware)
		}
	}
}
