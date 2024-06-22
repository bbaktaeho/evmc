package evmc

type debugNamespace struct {
	c caller
}

func (d *debugNamespace) TraceBlockByNumber() {}
