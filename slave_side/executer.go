package slave_side

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"vcbb/log"
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
	url    string
	logger *log.LoggerInstance
}

func NewPyExecuter(url string, logSystem *log.LogSystem) *PyExecuter {
	logger := logSystem.GetInstance(log.Topic("Slave Executer"))
	logger.Log("New Python Executer")
	return &PyExecuter{
		url:    url,
		logger: logger,
	}
}

func (this *PyExecuter) Run(partitionCnt, partitionIdOffset uint64, keys [][]string, code string, result chan *executeResult) {
	this.logger.Log("Try To Execute Code " + code)
	reqobj := executeReq{
		PartitionCnt:      partitionCnt,
		PartitionIdOffset: partitionIdOffset,
		Keys:              keys,
		Code:              code,
	}
	req, _ := json.Marshal(reqobj)
	res, err := http.Post(this.url+"/execute", "application/json", bytes.NewReader(req))
	if err != nil {
		this.logger.Err("Fail To Execute " + err.Error())
		result <- &executeResult{err: err}
		return
	}
	defer res.Body.Close()
	//fmt.Println(res)
	if res.Status != "200 OK" {
		this.logger.Err("Fail To Execute " + res.Status)
		result <- &executeResult{err: fmt.Errorf(res.Status)}
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		this.logger.Err("Fail To Read Execute Result " + err.Error())
		result <- &executeResult{err: err}
		return
	}
	var ans [][]string
	json.Unmarshal(body, &ans)
	this.logger.Log("Execute OK")
	result <- &executeResult{
		result: ans,
		err:    nil,
	}
}
