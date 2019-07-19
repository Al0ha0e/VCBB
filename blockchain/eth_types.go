package blockchain

type sendTransactionRes struct {
	Id      uint64 `json:"id"`
	JsonRPC string `json:"jsonrpc"`
	Tx      string `json:"result"`
}
