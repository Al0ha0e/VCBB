package types

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Address [20]byte

func (this Address) ToString() string {
	return hexutil.Encode(this[:])
}
func NewAddress(str string) Address {
	tmp, _ := hex.DecodeString(str[2:])
	var ret Address
	for i := 0; i < 20; i++ {
		ret[i] = tmp[i]
	}
	return ret
}

type Account struct {
	Id         Address
	PrivateKey string
}

func NewAccount(id [20]byte, privateKey string) *Account {
	return &Account{
		Id:         id,
		PrivateKey: privateKey,
	}
}
