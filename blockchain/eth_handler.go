package blockchain

import (
	"encoding/json"
	"vcbb/types"
)

type EthHandler struct {
	url string
}

func NewEthHandler(url string) *EthHandler {
	return &EthHandler{url: url}
}

func (this *EthHandler) GetTransaction(tx string) (interface{}, error) {
	req, err := genRPCReqJSON("2.0", "eth_getTransactionByHash", []string{tx}, "1")
	if err != nil {
		return nil, err
	}
	res, err := post(this.url, req)
	if err !=nil{
		return err
	}
	return ret,nil
}

func (this *EthHandler) SendTransaction(from, to types.Address, gas, gasPrice, value, data string) (string, error) {
	params := []string{from.ToString(), to.ToString(), gas, gasPrice, value, data}
	req, err := genRPCReqJSON("2.0", "eth_sendTransaction", params, "1")
	if err != nil {
		return "", err
	}
	res, err := post(this.url, req)
	if err != nil {
		return "", err
	}
	var txres sendTransactionRes
	err = json.Unmarshal([]byte(res), txres)
	return txres.Tx, nil
}

func (this *EthHandler) SendContract(from types.Address, gas, gasPrice, value, content string) (types.Address, error) {
	var to types.Address
	tx, err := this.SendTransaction(from, to, gas, gasPrice, value, content)

} /*
params: [{
  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
  "to": "0xd46e8dd67c5d32be8058bb8eb970870f072445675",
  "gas": "0x76c0", // 30400,
  "gasPrice": "0x9184e72a000", // 10000000000000
  "value": "0x9184e72a", // 2441406250
  "data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
}]*/
func (this *EthHandler) SendComputationContract(from types.Address, gas, gasPrice, value string) (types.Address, error) {

	return this.SendContract(from, gas, gasPrice, value)
}
