package net

type NetSimulator struct {
	userList map[string]chan []byte
}

func NewNetSimulator() *NetSimulator {
	return &NetSimulator{userList: make(map[string]chan []byte)}
}

func (this *NetSimulator) RegisterUser(name string, ch chan []byte) {
	this.userList[name] = ch
}

func (this *NetSimulator) SendMessageTo(name string, content []byte) {
	this.userList[name] <- content
}
