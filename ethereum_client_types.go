package ethereum_parser

import "encoding/json"

type requestBody struct {
	JsonRPC        string        `json:"jsonrpc"`
	EthereumMethod string        `json:"method"`
	Params         []interface{} `json:"params"`
	ID             int           `json:"id"`
}

type responseBody struct {
	ID      int             `json:"ID"`
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
}

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
}

type Block struct {
	Hash         string        `json:"hash"`
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

const (
	blockNumber       = "eth_blockNumber"
	getBlocksByNumber = "eth_getBlockByNumber"

	// ID Not sure about the ID, so I have left it as a const here
	ID = 1
)
