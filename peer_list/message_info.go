package peer_list

import "github.com/Al0ha0e/vcbb/types"

type MessageInfo struct {
	From        types.Address
	FromSession string
	Content     []byte
}

func NewMessageInfo(fr types.Address, frsess string, content []byte) MessageInfo {
	return MessageInfo{
		From:        fr,
		FromSession: frsess,
		Content:     content,
	}
}

type Message struct {
	From        types.Address `json:"from"`
	To          types.Address `json:"to"`
	FromSession string        `json:"fromSession"`
	ToSession   string        `json:"toSession"`
	Method      string        `json:"method"`
	Content     []byte        `json:"content"`
	Dist        uint8         `json:"dist"`
}

func newMessage(fr, to types.Address, frsess, tosess, method string, content []byte, dist uint8) Message {
	return Message{
		From:        fr,
		To:          to,
		FromSession: frsess,
		ToSession:   tosess,
		Method:      method,
		Content:     content,
		Dist:        dist,
	}
}
