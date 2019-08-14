package slave_side

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type executeReq struct {
	PartitionCnt uint64   `json:"partitionCnt"`
	Keys         []string `json:"keys"`
	Code         string   `json:"code"`
}

type executeResult struct {
	result [][]string
	err    error
}

type Executer interface {
	Run(partitionCnt uint64, keys []string, code string, result chan *executeResult)
}

type PyExecuter struct {
	url string
}

func NewPyExecuter(url string) *PyExecuter {
	return &PyExecuter{
		url: url,
	}
}

func (this *PyExecuter) Run(partitionCnt uint64, keys []string, code string, result chan *executeResult) {
	reqobj := executeReq{
		PartitionCnt: partitionCnt,
		Keys:         keys,
		Code:         code,
	}
	req, _ := json.Marshal(reqobj)
	res, err := http.Post(this.url, "application/json", bytes.NewReader(req))
	if err != nil {
		result <- &executeResult{err: err}
		return
	}
	defer res.Body.Close()
	//fmt.Println(res)
	if res.Status != "200 OK" {
		result <- &executeResult{err: fmt.Errorf(res.Status)}
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		result <- &executeResult{err: err}
		return
	}
	var ans [][]string
	json.Unmarshal(body, &ans)
	//fmt.Println(body, sb)
	result <- &executeResult{
		result: ans,
		err:    nil,
	}
}
