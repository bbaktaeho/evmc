package evmc

const (
	defaultConnPool int = 10
	// defaultIdleConnTimeout time.Duration = 30 * time.Second

	defaultMaxBatchItems int = 100
	defaultMaxBatchSize  int = 30 * 1024 * 1024
)

type options struct {
	connPool      int
	maxBatchItems int
	maxBatchSize  int
}

func newOps() *options {
	return &options{
		connPool:      defaultConnPool,
		maxBatchItems: defaultMaxBatchItems,
		maxBatchSize:  defaultMaxBatchSize,
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
