package ethereum_parser_test

import (
	"context"
	eth "ethereum_parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Implementing table-testing as well as simple test case to showcase different types ( test suit on api_test.go)
func TestEthereumClient_GetCurrentBlock(t *testing.T) {
	client := eth.NewEthereumClient(eth.EthereumClientConfig{Addr: addr, JsonRPC: ver})

	block, err := client.GetCurrentBlock(context.Background())
	assert.NoError(t, err)
	assert.Greater(t, block, int64(0))
}

func TestEthereumClient_GetBlockByNumber(t *testing.T) {
	client := eth.NewEthereumClient(eth.EthereumClientConfig{Addr: addr, JsonRPC: ver})

	tt := []struct {
		block        int64
		expectedHash string
	}{
		{
			block:        17085880,
			expectedHash: "0x08f2a2cd8370df0e926d3fdb92f4ebf61bf1c849f5b5967155eb82657b31b539",
		},
	}

	for _, testCase := range tt {
		receivedBlock, err := client.GetBlockByNumber(context.Background(), testCase.block)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedHash, receivedBlock.Hash)
	}
}

const addr = "https://cloudflare-eth.com"
const ver = "2.0"
