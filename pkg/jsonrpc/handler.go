package jsonrpc

// Handler should return either of Result or Error.
type Handler func(req *Request) (Result, *Error)
