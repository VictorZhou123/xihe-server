package main

import	(
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/points/messagequeue"
)
 

func (h *handler) HandleAddUserPoints(msg *comsg.MsgNormal) error {
	return messagequeue.Subscribe(h.point, []string{})
}