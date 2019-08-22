package net

import (
	"fmt"
	"net"
	"sync"
	"vcbb/types"
)

type UDPInfo struct {
	strAddr string
	udpAddr *net.UDPAddr
}

type UDPNetService struct {
	info     [2]*UDPInfo
	peers    map[string]*UDPInfo
	bus      chan []byte
	conn     *net.UDPConn
	buffSize uint64
	lock     sync.Mutex
}

func NewUDPInfo(addr string) (*UDPInfo, error) {
	udpaddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	return &UDPInfo{
		strAddr: addr,
		udpAddr: udpaddr,
	}, nil
}

func NewUDPNetService(inaddr, outaddr string) (*UDPNetService, error) {
	inInfo, err := NewUDPInfo(inaddr)
	if err != nil {
		return nil, err
	}
	outInfo, err := NewUDPInfo(outaddr)
	if err != nil {
		return nil, err
	}
	return &UDPNetService{
		info:     [2]*UDPInfo{inInfo, outInfo},
		peers:    make(map[string]*UDPInfo),
		buffSize: 4096,
	}, nil
}

func (this *UDPNetService) Run() error {
	var err error
	this.conn, err = net.ListenUDP("udp", this.info[0].udpAddr)
	if err != nil {
		return err
	}
	go this.serve()
	return nil
}

func (this *UDPNetService) serve() {
	for {
		buffer := make([]byte, this.buffSize)
		l, _, err := this.conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		this.bus <- buffer[:l]
	}
}

func (this *UDPNetService) RegisterUser(name string, ch chan []byte) {
	this.bus = ch
}
func (this *UDPNetService) SendMessageTo(name string, content []byte) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	info := this.peers[name]
	if info == nil {
		return fmt.Errorf("PEER NOT EXIST")
	}
	conn, err := net.DialUDP("udp", this.info[1].udpAddr, info.udpAddr)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(content, info.udpAddr)
	return err
}

func (this *UDPNetService) AddPeer(account types.Address, udpAddr string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	info, _ := NewUDPInfo(udpAddr)
	this.peers[account.ToString()] = info
}
