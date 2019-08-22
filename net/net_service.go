package net

type NetService interface {
	RegisterUser(name string, ch chan []byte)
	SendMessageTo(name string, content []byte)
}
