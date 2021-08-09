package jsonrpc

import log "github.com/sirupsen/logrus"

func PanicMiddleware(h Handler) Handler {
	return func(req *Request) (res Result, err *Error) {
		defer func() {
			if r := recover(); r != nil {
				err = &Error{
					Message: "panic occurred",
				}

				log.WithField("r", r).Error("panic caught")
			}
		}()

		return h(req)
	}
}
