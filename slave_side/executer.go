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

type Executer interface {
	Run()
}

type PyExecuter struct {
	code string
	url  string
}

func (this *PyExecuter) Run(partitionCnt uint64, keys []string, code string, result chan [][]string) {
	this.url = "http://127.0.0.1:8080/hello/test"
	reqobj := executeReq{
		PartitionCnt: partitionCnt,
		Keys:         keys,
		Code:         code,
	}
	req, _ := json.Marshal(reqobj)
	res, _ := http.Post(this.url, "application/json", bytes.NewReader(req))
	//fmt.Println(res)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var sb [][]string
	json.Unmarshal(body, &sb)
	fmt.Println(body, sb)
}
