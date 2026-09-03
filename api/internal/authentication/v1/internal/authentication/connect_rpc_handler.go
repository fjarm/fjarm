package authentication

import (
	"log/slog"

	"buf.build/go/protovalidate"
)

type AuthenticationServiceConnectRPCHandler struct {
	logger    *slog.Logger
	validator protovalidate.Validator
}
