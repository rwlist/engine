package jsonrpc

type Middleware func(Handler) Handler

func ApplyMiddlewares(h Handler, mdlw []Middleware) Handler {
	for i := 0; i < len(mdlw); i++ {
		h = mdlw[len(mdlw)-1-i](h)
	}

	return h
}
