package jsonrpc

type Middleware func(Handler) Handler
