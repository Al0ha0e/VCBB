package peer_list

import "github.com/Al0ha0e/vcbb/types"

type MessageInfo struct {
	From    types.Address
	Content []byte
}

type Message struct {
	From    types.Address `json:"from"`
	To      types.Address `json:"to"`
	Session string        `json:"session"`
	Method  string        `json:"method"`
	Content []byte        `json:"content"`
	Dist    uint8         `json:"dist"`
}

func newMessage(fr, to types.Address, sess, method string, content []byte, dist uint8) Message {
	return Message{
		From:    fr,
		To:      to,
		Session: sess,
		Method:  method,
		Content: content,
		Dist:    dist,
	}
}
