package types

type Address [20]byte

func (this Address) ToString() string {
	return string(this[:])
}

type Account struct {
	id Address
}

func NewAccount(id [20]byte) *Account {
	return &Account{id: id}
}
