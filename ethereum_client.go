package ethereum_parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var _ ethereumClient = EthereumClient{}

type EthereumClient struct {
	client  *http.Client
	rootUrl string
	jsonRPC string
}

type EthereumClientConfig struct {
	Addr    string `env:"ETHEREUM_CLIENT_URL" envDefault:"https://cloudflare-eth.com"`
	JsonRPC string `env:"JSON_RPC" envDefault:"2.0"`
}

func (c EthereumClient) call(ctx context.Context, method string, params []interface{}, v interface{}) error {
	body, err := json.Marshal(requestBody{JsonRPC: c.jsonRPC, EthereumMethod: method, Params: params, ID: ID})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rootUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK:
		var respBody responseBody
		// unmarshalling response and result from ethereum api
		if err = json.Unmarshal(bodyBytes, &respBody); err != nil {
			return fmt.Errorf("failed to unmarshal the response body: %v", err)
		}

		if err = json.Unmarshal([]byte(respBody.Result), &v); err != nil {
			return fmt.Errorf("failed to unmarshal the response body: %v", err)
		}

		return nil
	default:
		return fmt.Errorf("there was an error with the request")
	}
}

func (c EthereumClient) GetCurrentBlock(ctx context.Context) (int64, error) {
	var currentBlock string
	err := c.call(ctx, blockNumber, []interface{}{}, &currentBlock)
	if err != nil {
		return 0, err
	}
	quantity, err := hexDecoder(currentBlock)
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (c EthereumClient) GetBlockByNumber(ctx context.Context, number int64) (Block, error) {
	var result Block
	err := c.call(ctx, getBlocksByNumber, []interface{}{hexEndcoder(number), true}, &result)
	if err != nil {
		return Block{}, err
	}
	return result, nil
}

func NewEthereumClient(config EthereumClientConfig) EthereumClient {
	return EthereumClient{
		client:  http.DefaultClient,
		rootUrl: config.Addr,
		jsonRPC: config.JsonRPC,
	}
}

func hexDecoder(hexValue string) (int64, error) {
	intValue, err := strconv.ParseInt(hexValue, 0, 64)
	if err != nil {
		fmt.Println("Error decoding hex:", err)
		return 0, err
	}

	return intValue, nil
}

func hexEndcoder(decValue int64) string {
	return fmt.Sprintf("0x%X", decValue)
}

type ethereumClient interface {
	// GetCurrentBlock retrieves the number of the most recent block
	GetCurrentBlock(ctx context.Context) (int64, error)

	// GetBlockByNumber returns information about a block by block number.
	GetBlockByNumber(ctx context.Context, number int64) (Block, error)
}
