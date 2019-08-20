package types

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
	Id              Address
	PrivateKey      string
	ECDSAPrivateKey *ecdsa.PrivateKey
}

func NewAccount(id [20]byte, privateKey string) *Account {
	ecdp, _ := crypto.HexToECDSA(privateKey)
	return &Account{
		Id:              id,
		PrivateKey:      privateKey,
		ECDSAPrivateKey: ecdp,
	}
}
