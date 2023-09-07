package messages

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	competitionmsg "github.com/opensourceways/xihe-server/competition/domain/message"
)

func (s sender) NotifyCalcScore(v *domain.SubmissionMessage) error {
	return s.send(topics.Submission, v)
}

func (s sender) SendCompetitionMsg(v *competitionmsg.MsgTask) error {
	return s.send(topics.Competition, v)
}
