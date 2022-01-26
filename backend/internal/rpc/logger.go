package rpc

import (
	"github.com/rwlist/gjrpc/pkg/jsonrpc"
	log "github.com/sirupsen/logrus"
)

func LogMiddleware() jsonrpc.Middleware {
	return func(next jsonrpc.Handler) jsonrpc.Handler {
		return func(req *jsonrpc.Request) (jsonrpc.Result, *jsonrpc.Error) {
			res, err := next(req)
			log.
				WithField("params", req.Params).
				WithField("result", res).
				WithField("rpc_error", err).
				Info("rpc call finished")
			return res, err
		}
	}
}
