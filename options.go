package evmc

import "time"

const (
	defaultConnPool        int           = 10
	defaultReqTimeout      time.Duration = 1 * time.Minute
	defaultIdleConnTimeout time.Duration = 2 * time.Minute

	defaultMaxBatchItems    int = 100
	defaultMaxBatchSize     int = 30 * 1024 * 1024
	defaultBatchCallWorkers int = 3

	defaultWsReadBufferSize   int = 1024
	defaultWsWriteBufferSize  int = 1024
	defaultWsMessageSizeLimit int = 0 // unlimited
	// defaultWsPingInterval     time.Duration = 30 * time.Second
	// defaultWsPingWriteTimeout time.Duration = 2 * time.Second
	// defaultWsPongTimeout      time.Duration = 10 * time.Second
)

type options struct {
	connPool         int
	reqTimeout       time.Duration
	idleConnTimeout  time.Duration
	maxBatchItems    int
	maxBatchSize     int
	batchCallWorkers int

	wsReadBufferSize   int
	wsWriteBufferSize  int
	wsMessageSizeLimit int
	// wsPingInterval     time.Duration
	// wsPingWriteTimeout time.Duration
	// wsPongTimeout      time.Duration
}

func newOps() *options {
	return &options{
		connPool:           defaultConnPool,
		reqTimeout:         defaultReqTimeout,
		idleConnTimeout:    defaultIdleConnTimeout,
		maxBatchItems:      defaultMaxBatchItems,
		maxBatchSize:       defaultMaxBatchSize,
		batchCallWorkers:   defaultBatchCallWorkers,
		wsReadBufferSize:   defaultWsReadBufferSize,
		wsWriteBufferSize:  defaultWsWriteBufferSize,
		wsMessageSizeLimit: defaultWsMessageSizeLimit,
		// wsPingInterval:     defaultWsPingInterval,
		// wsPingWriteTimeout: defaultWsPingWriteTimeout,
		// wsPongTimeout:      defaultWsPongTimeout,
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
		o.maxBatchItems = items
	})
}

func WithMaxBatchSize(size int) Options {
	return optionFunc(func(o *options) {
		o.maxBatchSize = size
	})
}

func WithBatchCallWorkers(workers int) Options {
	if workers < 1 {
		workers = defaultBatchCallWorkers
	}
	return optionFunc(func(o *options) {
		o.batchCallWorkers = workers
	})
}

func WithWsReadBufferSize(size int) Options {
	return optionFunc(func(o *options) {
		o.wsReadBufferSize = size
	})
}

func WithWsWriteBufferSize(size int) Options {
	return optionFunc(func(o *options) {
		o.wsWriteBufferSize = size
	})
}

func WithWsMessageSizeLimit(limit int) Options {
	return optionFunc(func(o *options) {
		o.wsMessageSizeLimit = limit
	})
}

// func WithWsPingInterval(interval time.Duration) Options {
// 	return optionFunc(func(o *options) {
// 		o.wsPingInterval = interval
// 	})
// }

// func WithWsPingWriteTimeout(timeout time.Duration) Options {
// 	return optionFunc(func(o *options) {
// 		o.wsPingWriteTimeout = timeout
// 	})
// }

// func WithWsPongTimeout(timeout time.Duration) Options {
// 	return optionFunc(func(o *options) {
// 		o.wsPongTimeout = timeout
// 	})
// }
