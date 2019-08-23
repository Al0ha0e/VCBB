package slave_side

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type executeReq struct {
	PartitionCnt      uint64     `json:"partitionCnt"`
	PartitionIdOffset uint64     `json:"partitionIdOffset"`
	Keys              [][]string `json:"keys"`
	Code              string     `json:"code"`
}

type executeResult struct {
	result [][]string
	err    error
}

type Executer interface {
	Run(partitionCnt, partitionIdOffset uint64, keys [][]string, code string, result chan *executeResult)
}

type PyExecuter struct {
	url string
}

func NewPyExecuter(url string) *PyExecuter {
	return &PyExecuter{
		url: url,
	}
}

func (this *PyExecuter) Run(partitionCnt, partitionIdOffset uint64, keys [][]string, code string, result chan *executeResult) {
	reqobj := executeReq{
		PartitionCnt:      partitionCnt,
		PartitionIdOffset: partitionIdOffset,
		Keys:              keys,
		Code:              code,
	}
	req, _ := json.Marshal(reqobj)
	res, err := http.Post(this.url+"/execute", "application/json", bytes.NewReader(req))
	if err != nil {
		fmt.Println("EXECUTER ERR", err)
		result <- &executeResult{err: err}
		return
	}
	defer res.Body.Close()
	//fmt.Println(res)
	if res.Status != "200 OK" {
		fmt.Println("EXECUTER ERR", res.Status)
		result <- &executeResult{err: fmt.Errorf(res.Status)}
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("EXECUTER ERR", err)
		result <- &executeResult{err: err}
		return
	}
	var ans [][]string
	json.Unmarshal(body, &ans)
	//fmt.Println(body, sb)
	fmt.Println("EXECUTER ANS", ans)
	result <- &executeResult{
		result: ans,
		err:    nil,
	}
}
