package net

import (
	"fmt"
	"testing"
	"vcbb/types"
)

func watch(addr string, ch chan []byte) {
	for {
		msg := <-ch
		fmt.Println(addr, msg)
	}
}

func TestUDPService(t *testing.T) {
	addr1 := "127.0.0.1:8081"
	addr2 := "127.0.0.1:8082"
	addr3 := "127.0.0.1:8083"
	addr4 := "127.0.0.1:8084"
	service1, err := NewUDPNetService(addr1, addr2)
	if err != nil {
		t.Error(err)
		return
	}
	service2, err := NewUDPNetService(addr3, addr4)
	if err != nil {
		t.Error(err)
		return
	}
	service1.Run()
	service2.Run()
	ch1 := make(chan []byte, 1)
	ch2 := make(chan []byte, 1)
	service1.RegisterUser("", ch1)
	service2.RegisterUser("", ch2)
	go watch(addr1, ch1)
	go watch(addr3, ch2)
	service1.AddPeer(types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d2"), "127.0.0.1:8083")
	service2.AddPeer(types.NewAddress("0xd247126aa720779a4172b88405ec2ad29c6a22d1"), "127.0.0.1:8081")
	service1.SendMessageTo("0xd247126aa720779a4172b88405ec2ad29c6a22d2", []byte{1, 2, 3})
	service2.SendMessageTo("0xd247126aa720779a4172b88405ec2ad29c6a22d1", []byte{4, 5, 6})
	for {

	}
}
