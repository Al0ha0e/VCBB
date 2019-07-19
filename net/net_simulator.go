package net

import (
	"vcbb/event"
)

const (
	SimuNetMsg = "SimulatorNetMessage"
)

type NetSimulator struct {
	eventSystem *event.EventSystem
	userList    map[string]chan string
}

func NewNetSimulator(eventSystem *event.EventSystem) *NetSimulator {
	return &NetSimulator{eventSystem: eventSystem}
}

func (this *NetSimulator) RegisterUser(name string, ch chan string) {
	this.userList[name] = ch
}

func (this *NetSimulator) BroadCast(content string) {
	this.eventSystem.Emit(SimuNetMsg, content)
}

func (this *NetSimulator) SendMessageTo(name string, content string) {
	this.userList[name] <- content
}
