package evmc

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func testEvmcForDebug() debugNamespace {
	rpcURL := "https://ethereum-mainnet.nodit.io/<api-key>"
	client, err := New(rpcURL)
	if err != nil {
		panic(err)
	}
	return debugNamespace{c: client}
}

// FIXME: internal JSON-RPC error
func Test_debugNamespace_TraceBlockByNumber(t *testing.T) {
	chain := testEvmcForDebug()
	traceResult, err := chain.TraceBlockByNumber(20000000, &TraceConfig{
		DefaultTraceConfig: &DefaultTraceConfig{
			EnableMemory:     true,
			DisableStack:     true,
			DisableStorage:   true,
			EnableReturnData: true,
			Debug:            true,
		},
		Timeout: time.Minute.String(),
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(traceResult, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debug_traceBlockByNumber.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}

func Test_debugNamespace_TraceBlockByNumber_4byteTracer(t *testing.T) {
	chain := testEvmcForDebug()
	traceResult, err := chain.TraceBlockByNumber(20000000, &TraceConfig{Tracer: FourByteTracer})
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(traceResult, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debug_traceBlockByNumber_4byteTracer.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}

func Test_debugNamespace_TraceBlockByNumber_muxTracer(t *testing.T) {
	chain := testEvmcForDebug()
	traceResult, err := chain.TraceBlockByNumber(20000000, &TraceConfig{Tracer: MuxTracer})
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(traceResult, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debug_traceBlockByNumber_muxTracer.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}

func Test_debugNamespace_TraceBlockByNumber_prestateTracer(t *testing.T) {
	chain := testEvmcForDebug()
	traceResult, err := chain.TraceBlockByNumber(20000000, &TraceConfig{Tracer: PrestateTracer})
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(traceResult, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debug_traceBlockByNumber_prestateTracer.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}

func Test_debugNamespace_TraceBlockByNumber_callTracer(t *testing.T) {
	chain := testEvmcForDebug()
	callTracer, err := chain.TraceBlockByNumber_callTracer(20000000, time.Second*100, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(callTracer, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debug_traceBlockByNumber_callTracer.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()

	callTracer, err = chain.TraceBlockByNumber_callTracer(20000000, time.Second*100, nil, &CallTracerConfig{
		WithLog: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err = json.MarshalIndent(callTracer, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ = os.OpenFile("debug_traceBlockByNumber_callTracer_withLog.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}

func Test_debugNamespace_TraceBlockByNumber_flatCallTracer(t *testing.T) {
	chain := testEvmcForDebug()
	callTracer, err := chain.TraceBlockByNumber_flatCallTracer(133294839, time.Second*100, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(callTracer, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debug_traceBlockByNumber_flatCallTracer.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()

	callTracer, err = chain.TraceBlockByNumber_flatCallTracer(133294839, time.Second*100, nil, &FlatCallTracerConfig{
		ConvertParityErrors: true,
		IncludePrecompiles:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	b, err = json.MarshalIndent(callTracer, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ = os.OpenFile("debug_traceBlockByNumber_flatCallTracer_convertParityErrors_includePrecompiles.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}

func Test_debugNamespace_TraceTransaction_callTracer(t *testing.T) {
	chain := testEvmcForDebug()
	callTracer, err := chain.TraceTransaction_callTracer("0x92826a6a3c5bee87f0a834d095817a41e426e815506acbe9e6174160bdb745d4", time.Second*100, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(callTracer, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	f, _ := os.OpenFile("debugTraceTransaction_callTracer.json", os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(string(b))
	f.Close()
}
