package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RPCReq struct {
	Version string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      string   `json:"id"`
}

func genRPCReqJSON(version string, method string, params []string, id string) (string, error) {
	retobj := RPCReq{
		Version: version,
		Method:  method,
		Params:  params,
		Id:      id,
	}
	ret, err := json.Marshal(retobj)
	return string(ret), err
}

//"http://127.0.0.1:8545"
func post(url string, msg string) (string, error) {
	res, err := http.Post(url, "application/json", strings.NewReader(msg))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if res.Status == "200 OK" {
		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		return string(body), nil
	} else {
		fmt.Println(res.Status)
		return "FAIL", nil
	}
	//return "FAIL", nil
}
