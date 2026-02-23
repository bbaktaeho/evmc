package evmc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// mockRPCServer는 테스트용 JSON-RPC HTTP 서버.
type mockRPCServer struct {
	mu       sync.Mutex
	handlers map[string]func(params json.RawMessage) interface{}
	server   *httptest.Server
}

func newMockRPCServer(t *testing.T) *mockRPCServer {
	t.Helper()
	m := &mockRPCServer{
		handlers: make(map[string]func(params json.RawMessage) interface{}),
	}
	m.server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.server.Close)
	return m
}

// writeJSON은 응답 헤더를 설정하고 v를 JSON으로 인코딩해 w에 쓴다.
// 인코딩 중 에러가 발생하면 500 Internal Server Error를 반환한다.
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, fmt.Sprintf("json encode error: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (m *mockRPCServer) handle(w http.ResponseWriter, r *http.Request) {
	var body json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// batch 요청 처리
	if len(body) > 0 && body[0] == '[' {
		var reqs []struct {
			ID     interface{}     `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal(body, &reqs); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		responses := make([]map[string]interface{}, 0, len(reqs))
		for _, req := range reqs {
			responses = append(responses, m.dispatch(req.ID, req.Method, req.Params))
		}
		writeJSON(w, responses)
		return
	}

	// 단일 요청 처리
	var req struct {
		ID     interface{}     `json:"id"`
		Method string          `json:"method"`
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	writeJSON(w, m.dispatch(req.ID, req.Method, req.Params))
}

func (m *mockRPCServer) dispatch(id interface{}, method string, params json.RawMessage) map[string]interface{} {
	m.mu.Lock()
	handler, ok := m.handlers[method]
	m.mu.Unlock()

	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"error": map[string]interface{}{
				// JSON-RPC 2.0 error code: method not found
				"code":    -32601,
				"message": fmt.Sprintf("method not found: %s", method),
			},
		}
	}
	result := handler(params)
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
}

func (m *mockRPCServer) on(method string, handler func(params json.RawMessage) interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[method] = handler
}

func (m *mockRPCServer) url() string {
	return m.server.URL
}
