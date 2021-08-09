package rwserv

import "github.com/rwlist/engine/pkg/jsonrpc"

type Server struct {
}

func (s *Server) Handle(req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
	return true, nil
}
