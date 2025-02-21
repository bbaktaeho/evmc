package evmc

import "errors"

var (
	ErrPendingBlockNotSupported           = errors.New("pending block is not supported")
	ErrWebsocketRequired                  = errors.New("websocket is required for subscriptions")
	ErrWalletRequired                     = errors.New("wallet is required")
	ErrTxGasLimitZero                     = errors.New("gas limit is zero")
	ErrTxGasPriceLessThanZero             = errors.New("gas price less than zero")
	ErrTxValueLessThanZero                = errors.New("value less than zero")
	ErrTxMaxFeePerGasLessThanZero         = errors.New("max fee per gas less than zero")
	ErrTxMaxPriorityFeePerGasLessThanZero = errors.New("max priority fee per gas less than zero")
	ErrFromRequired                       = errors.New("from address is required")
	ErrToRequired                         = errors.New("to address is required")
	ErrTxRequired                         = errors.New("tx is required")
	ErrInvalidRange                       = errors.New("invalid range from > to")
	ErrChainIDLessThanZero                = errors.New("chain id is required")
)
