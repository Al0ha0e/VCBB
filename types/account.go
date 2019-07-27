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
	id Address
}

func NewAccount(id [20]byte) *Account {
	return &Account{id: id}
}
