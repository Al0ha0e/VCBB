package net

import "vcbb/types"

type NetService interface {
	RegisterUser(name string, ch chan []byte)
	SendMessageTo(name string, content []byte)
	Run() error
	AddPeer(account types.Address, udpAddr string)
}
