package evmc

import "errors"

var (
	ErrPendingBlockNotSupported = errors.New("pending block is not supported")
)
