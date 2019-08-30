package log

import (
	"fmt"
	"os"
	"sync"
)

type Topic string

type LogSystem struct {
	path      string
	file      *os.File
	instances map[Topic]*LoggerInstance
	lock      sync.Mutex
	logLock   sync.Mutex
	//logger    func(string) error
}

type LoggerInstance struct {
	logSystem      *LogSystem
	LogTopic       Topic
	subInstances   map[Topic]*LoggerInstance
	fatherInstance *LoggerInstance
}

func NewLogSystem(path string) (*LogSystem, error) {
	var file *os.File
	var err error
	if len(path) > 0 {
		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0766)
		if err != nil {
			return nil, err
		}
	}
	return &LogSystem{
		path:      path,
		file:      file,
		instances: make(map[Topic]*LoggerInstance),
		//logger:    logger,
	}, nil
}

func newLoggerInstance(topic Topic) *LoggerInstance {
	return &LoggerInstance{
		LogTopic:     topic,
		subInstances: make(map[Topic]*LoggerInstance),
	}
}

func (this *LogSystem) GetInstance(topic Topic) *LoggerInstance {
	this.lock.Lock()
	defer this.lock.Unlock()
	inst := this.instances[topic]
	if inst != nil {
		return inst
	}
	inst = newLoggerInstance(topic)
	inst.logSystem = this
	this.instances[topic] = inst
	return inst
}

func (this *LoggerInstance) GetSubInstance(topic Topic) *LoggerInstance {
	inst := this.subInstances[topic]
	if inst != nil {
		return inst
	}
	inst = newLoggerInstance(topic)
	inst.fatherInstance = this
	this.subInstances[topic] = inst
	return inst
}

func (this *LogSystem) Log(msg string) error {
	fmt.Println(msg)
	return nil
	/*
		if this.file == nil {
			return nil
		}
		this.logLock.Lock()
		defer this.logLock.Unlock()
		_, err := this.file.WriteString(msg)
		return err*/
}

func (this *LogSystem) Err(msg string) error {
	fmt.Println("Error:", msg)
	return nil
}

/*
func WriteString(msg string) error {

}*/

func (this *LoggerInstance) Log(msg string) error {
	if this.fatherInstance == nil {
		return this.logSystem.Log(string(this.LogTopic) + " " + msg /* + "\r\n"*/)
	} else {
		return this.fatherInstance.Log(string(this.LogTopic) + " " + msg)
	}
}

func (this *LoggerInstance) Err(msg string) error {
	if this.fatherInstance == nil {
		return this.logSystem.Err(string(this.LogTopic) + " " + msg /*+ "\r\n"*/)
	} else {
		return this.fatherInstance.Err(string(this.LogTopic) + " " + msg)
	}
}
