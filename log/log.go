package log

import (
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
}

type LoggerInstance struct {
	logSystem *LogSystem
	LogTopic  Topic
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
	}, nil
}

func newLoggerInstance(topic Topic) *LoggerInstance {
	return &LoggerInstance{
		LogTopic: topic,
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

func (this *LogSystem) writeString(msg string) error {
	if this.file == nil {
		return nil
	}
	this.logLock.Lock()
	defer this.logLock.Unlock()
	_, err := this.file.WriteString(msg)
	return err
}

func (this *LoggerInstance) Log(msg string) error {
	return this.logSystem.writeString(string(this.LogTopic) + " :" + msg + "\r\n")
}
