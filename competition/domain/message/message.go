package message

import (
	comsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/competition/domain"
)

type MsgTask comsg.MsgNormal

type CompetitionMessageProducer interface {
	NotifyCalcScore(*domain.SubmissionMessage) error
	SendCompetitionMsg(*MsgTask) error
}
