package ethereum_parser_test

import (
	"context"
	eth "ethereum_parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Implementing a simple test that retrieves the current block ( not a test suite this time) just to showcase use of different types
func TestCurrentBlock(t *testing.T) {
	client := eth.NewEthereumClient(eth.EthereumClientConfig{Addr: addr, JsonRPC: ver})

	block, err := client.GetCurrentBlock(context.Background(), eth.ID)
	assert.NoError(t, err)
	assert.Equal(t, block, int64(0))
}

const addr = "https://cloudflare-eth.com"
const ver = "2.0"
