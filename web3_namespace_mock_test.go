package evmc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_web3Namespace_mock_ClientVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "geth",
			version: "Geth/v1.13.0-stable-64fe2a41/linux-amd64/go1.21.0",
			want:    "Geth/v1.13.0-stable-64fe2a41/linux-amd64/go1.21.0",
		},
		{
			name:    "erigon",
			version: "erigon/v2.52.0/linux-amd64/go1.21.0",
			want:    "erigon/v2.52.0/linux-amd64/go1.21.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRPCServer(t)
			mock.on("web3_clientVersion", func(params json.RawMessage) interface{} {
				return tt.version
			})

			client := testEvmc(mock.url())
			version, err := client.Web3().ClientVersion()
			require.NoError(t, err)
			assert.Equal(t, tt.want, version)
		})
	}
}

func Test_evmc_NodeClient(t *testing.T) {
	mock := newMockRPCServer(t)
	mock.on("web3_clientVersion", func(params json.RawMessage) interface{} {
		return "Geth/v1.13.0-stable/linux-amd64/go1.21.0"
	})

	client := testEvmc(mock.url())
	name, version, err := client.NodeClient()
	require.NoError(t, err)
	assert.Equal(t, "Geth", name)
	assert.Equal(t, "v1.13.0-stable", version)
}

func Test_evmc_NodeClient_ShortVersion(t *testing.T) {
	mock := newMockRPCServer(t)
	mock.on("web3_clientVersion", func(params json.RawMessage) interface{} {
		return "SimpleClient"
	})

	client := testEvmc(mock.url())
	name, version, err := client.NodeClient()
	require.NoError(t, err)
	// 슬래시가 없으면 name, version이 빈 문자열
	assert.Empty(t, name)
	assert.Empty(t, version)
}
