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

// There is some repetition we can create a doer that executes the request and passes on the response body for further processing by each method.

func (c EthereumClient) GetCurrentBlock(ctx context.Context, id int) (int64, error) {
	reqBody, err := json.Marshal(requestBody{JsonRPC: c.jsonRPC, EthereumMethod: blockNumber, Params: []interface{}{}, ID: id})
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rootUrl, bytes.NewReader(reqBody))
	if err != nil {
		return 0, err
	}

	request.Header.Add("Content-Type", "application/json")

	// SENDING
	response, err := c.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var result responseBody
	switch response.StatusCode {
	case http.StatusOK:

		err := json.NewDecoder(response.Body).Decode(&result)
		if err = json.Unmarshal(respBody, &result); err != nil {
			return 0, fmt.Errorf("failed to umarshal response body: %v", err)
		}

		var currentBlock string
		if err = json.Unmarshal([]byte(result.Result), &currentBlock); err != nil {
			return 0, fmt.Errorf("failed to unmarshal response body: %v", err)

		}
		quantity, err := hexDecoder(currentBlock)
		if err != nil {
			return 0, err
		}
		return quantity, nil
	default:
		return 0, fmt.Errorf("there was an error with the resuqest")
	}
}

func (c EthereumClient) GetBlockByNumber(ctx context.Context, number int64) (Block, error) {
	body, err := json.Marshal(requestBody{JsonRPC: c.jsonRPC, EthereumMethod: getBlocksByNumber, Params: []interface{}{hexEndcoder(number), true}, ID: ID})
	if err != nil {
		return Block{}, fmt.Errorf("failed to marshal request body: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.rootUrl, bytes.NewReader(body))
	if err != nil {
		return Block{}, err
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return Block{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Block{}, err
	}

	var result Block
	switch response.StatusCode {
	case http.StatusOK:
		var respBody responseBody
		// unmarshalling response and result from ethereum api
		if err = json.Unmarshal(bodyBytes, &respBody); err != nil {
			return Block{}, fmt.Errorf("failed to unmarshaling the response body: %v", err)
		}

		if err = json.Unmarshal([]byte(respBody.Result), &result); err != nil {
			return Block{}, fmt.Errorf("failed to unmarshal the response body: %v", err)
		}

		return result, nil
	default:
		return Block{}, fmt.Errorf("there was an error with the resuqest")
	}
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
	GetCurrentBlock(ctx context.Context, id int) (int64, error)

	// GetBlockBlockByNumber GetBlockByNumber returns information about a block by block number.
	GetBlockByNumber(ctx context.Context, number int64) (Block, error)
}
