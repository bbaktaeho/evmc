package evmc

import "errors"

var (
	ErrPendingBlockNotSupported   = errors.New("pending block is not supported")
	ErrWebsocketRequired          = errors.New("websocket is required for subscriptions")
	ErrWalletRequired             = errors.New("wallet is required")
	ErrTxGasLimitZero             = errors.New("gas limit is zero")
	ErrTxGasPriceZero             = errors.New("gas price is zero")
	ErrTxValueLessThanZero        = errors.New("value less than zero")
	ErrTxMaxFeePerGasZero         = errors.New("max fee per gas is zero")
	ErrTxMaxPriorityFeePerGasZero = errors.New("max priority fee per gas is zero")
	ErrFromRequired               = errors.New("from address is required")
	ErrToRequired                 = errors.New("to address is required")
	ErrTxRequired                 = errors.New("tx is required")
)
