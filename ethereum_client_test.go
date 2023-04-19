package ethereum_parser

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEthereumClient(t *testing.T) {
	client := NewEthereumClient(EthereumClientConfig{})

	block, err := client.GetCurrentBlock(context.Background(), ID)
	assert.NoError(t, err)
	assert.Equal(t, block, int64(0))
}
