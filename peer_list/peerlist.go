package peer_list

import "vcbb/types"

type PeerList struct {
}

func (this *PeerList) BroadCast([]byte) error {
	return nil
}

func (this *PeerList) SendMsgTo(to types.Address, msg []byte) {

}

func (this *PeerList) UpdatePunishedPeers(map[string][]types.Address) {

}
