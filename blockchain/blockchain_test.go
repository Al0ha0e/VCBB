package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"vcbb/types"
)

func TestCalcContract(t *testing.T) {
	addr := types.NewAddress("0x2acaac851b020ceb644bc506a3a932f4d0867afd")
	pvi := "3b82b9641714c4bb9a3e3a23ca9e8170772fcdeedd9e4591e7d03ebe564a579e"
	acco := types.Account{
		Id:         addr,
		PrivateKey: pvi,
	}
	hdl, err := NewEthBlockChainHandler("http://127.0.0.1:8545", acco)
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
}
