package evmc

import "errors"

var (
	ErrPendingBlockNotSupported = errors.New("pending block is not supported")
	ErrWebsocketRequired        = errors.New("websocket is required for subscriptions")
)
