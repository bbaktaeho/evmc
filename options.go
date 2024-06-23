package evmc

import "time"

const (
	defaultConnPool        int           = 10
	defaultReqTimeout      time.Duration = 1 * time.Minute
	defaultIdleConnTimeout time.Duration = 2 * time.Minute

	defaultMaxBatchItems int = 100
	defaultMaxBatchSize  int = 30 * 1024 * 1024
)

type options struct {
	connPool        int
	reqTimeout      time.Duration
	idleConnTimeout time.Duration
	maxBatchItems   int
	maxBatchSize    int
}

func newOps() *options {
	return &options{
		connPool:        defaultConnPool,
		reqTimeout:      defaultReqTimeout,
		idleConnTimeout: defaultIdleConnTimeout,
		maxBatchItems:   defaultMaxBatchItems,
		maxBatchSize:    defaultMaxBatchSize,
	}
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

type Options interface {
	apply(*options)
}

func WithConnPool(pool int) Options {
	return optionFunc(func(o *options) {
		o.connPool = pool
	})
}

func WithReqTimeout(timeout time.Duration) Options {
	return optionFunc(func(o *options) {
		o.reqTimeout = timeout
	})
}

func WithMaxBatchItems(items int) Options {
	return optionFunc(func(o *options) {
		o.connPool = items
	})
}

func WithMaxBatchSize(size int) Options {
	return optionFunc(func(o *options) {
		o.connPool = size
	})
}
