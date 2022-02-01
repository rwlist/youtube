package proto

import "github.com/rwlist/gjrpc/pkg/jsonrpc"

var (
	InternalError = jsonrpc.Error{Code: 1, Message: "internal error"}
	AuthError     = jsonrpc.Error{Code: 2, Message: "auth is invalid or expired, reset access_key if any"}
)
