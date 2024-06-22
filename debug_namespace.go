package evmc

const (
	debugTraceBlockByNumber = "debug_traceBlockByNumber"
	debugTraceTransaction   = "debug_traceTransaction"
)

type debugNamespace struct {
	c caller
}

func (d *debugNamespace) TraceBlockByNumber() {}
