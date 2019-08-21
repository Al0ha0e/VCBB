package blockchain

import (
	"fmt"
	"math/big"
	"testing"
	"vcbb/types"
)

/*
func TestCalcContractDeploy(t *testing.T) {
	addr := types.NewAddress("0x2acaac851b020ceb644bc506a3a932f4d0867afd")
	pvi := "3b82b9641714c4bb9a3e3a23ca9e8170772fcdeedd9e4591e7d03ebe564a579e"
	acco := types.NewAccount(addr, pvi)
	hdl, err := NewEthBlockChainHandler("ws://127.0.0.1:8546", acco)
	if err != nil {
		t.Error("HANDLER", err)
	}
	fmt.Println(hdl.client)
	gp, err := hdl.client.SuggestGasPrice(context.Background())
	if err != nil {
		t.Error("GP", err)
	}
	binfo := &ContractDeployInfo{
		Value:    big.NewInt(130),
		GasPrice: gp,
		GasLimit: uint64(4712388),
	}
	cinfo := &CalculationContractDeployInfo{
		Id:               "test",
		St:               big.NewInt(0),
		Fund:             big.NewInt(100),
		ParticipantCount: uint8(2),
		Distribute:       [8]*big.Int{big.NewInt(20), big.NewInt(10), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
	}
	ct := NewCalculationContract(hdl, nil, binfo, cinfo)
	caddr, err := ct.Start()
	if err != nil {
		t.Error("DEPLOY", err)
	}
	fmt.Println(caddr)
}*/

func TestCalcContractCommit(t *testing.T) {
	addr := types.NewAddress("0x2acaac851b020ceb644bc506a3a932f4d0867afd")
	pvi := "3b82b9641714c4bb9a3e3a23ca9e8170772fcdeedd9e4591e7d03ebe564a579e"
	acco := types.NewAccount(addr, pvi)
	hdl, err := NewEthBlockChainHandler("ws://127.0.0.1:8546", acco)
	if err != nil {
		t.Error("HANDLER", err)
	}
	fmt.Println(hdl.client)
	binfo := &ContractDeployInfo{
		Value:    big.NewInt(130),
		GasLimit: uint64(4712388),
	}
	cinfo := &CalculationContractDeployInfo{
		Id:               "test",
		St:               big.NewInt(0),
		Fund:             big.NewInt(100),
		ParticipantCount: uint8(2),
		Distribute:       [8]*big.Int{big.NewInt(20), big.NewInt(10), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)},
	}
	up := make(chan *Answer, 5)
	ct := NewCalculationContract(hdl, up, binfo, cinfo)
	go func() {
		for {
			ans := <-up
			fmt.Println("NEW ANSWER", ans)
		}
	}()
	caddr, err := ct.Start()
	if err != nil {
		t.Error("DEPLOY", err)
	}
	fmt.Println(caddr)
	addr2 := types.NewAddress("0x9c67d6e615fb9fb28ddad773fbcfa8e5dad092f3")
	pri2 := "ee09c465edc1674d382157f9edb26681707b79b31cab452450776a2a1ad57be5"
	acco2 := types.NewAccount(addr2, pri2)
	hdl2, _ := NewEthBlockChainHandler("ws://localhost:8546", acco2)
	ct2, _ := CalculationContractFromAddress(hdl2, caddr)
	info2 := &ContractDeployInfo{
		Value:    big.NewInt(100),
		GasLimit: uint64(4712388),
	}
	sb := make([][]string, 2)
	sb[0] = []string{"a", "b"}
	sb[1] = []string{"c", "d"}
	err = ct2.Commit(info2, sb, "sb")
	fmt.Println("COMMIT OK")
	if err != nil {
		t.Error("COMMIT ERR", err)
	}
	for {
	}
}
