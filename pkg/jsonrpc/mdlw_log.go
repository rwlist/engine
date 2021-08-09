package jsonrpc

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type LogOptions struct {
	Logger      log.FieldLogger
	IncludeBody bool
}

func LogMiddleware(opts *LogOptions) func(handler Handler) Handler {
	if opts == nil {
		opts = &LogOptions{}
	}
	if opts.Logger == nil {
		opts.Logger = log.StandardLogger()
	}

	return func(handler Handler) Handler {
		return func(req *Request) (Result, *Error) {
			start := time.Now()

			res, err := handler(req)

			logger := opts.Logger.
				WithField("method", req.Method).
				WithField("duration", time.Since(start).String())

			if opts.IncludeBody {
				logger = logger.WithField("params", req.Params)

				if err != nil {
					logger = logger.WithField("error", err)
				} else {
					logger = logger.WithField("result", res)
				}
			}

			logger.Info("request finished")

			return res, err
		}
	}
}
