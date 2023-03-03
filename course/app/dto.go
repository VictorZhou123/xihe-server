package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/course/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

// player apply
type PlayerApplyCmd domain.Player

func (cmd *PlayerApplyCmd) Validate() error {
	b := cmd.Student.Account != nil &&
		cmd.Student.Name != nil &&
		cmd.Student.Email != nil &&
		cmd.Student.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *PlayerApplyCmd) toPlayer() (p domain.Player) {
	return *(*domain.Player)(cmd)
}

// list
type CourseListCmd struct {
	Status domain.CourseStatus
	Type   domain.CourseType
	User   types.Account
}

type CourseSummaryDTO struct {
	PlayerCount int    `json:"count"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Hours       int    `json:"hours"`
	Host        string `json:"host"`
	Desc        string `json:"desc"`
	Status      string `json:"status"`
	Poster      string `json:"poster"`
	Duration    string `json:"duration"`
	Type        string `json:"type"`
}

func (s courseService) toCourseSummaryDTO(
	c *domain.CourseSummary, playerCount int, dto *CourseSummaryDTO,
) {
	*dto = CourseSummaryDTO{
		PlayerCount: playerCount,
		Id:          c.Id,
		Name:        c.Name.CourseName(),
		Hours:       c.Hours.CourseHours(),
		Host:        c.Host.CourseHost(),
		Desc:        c.Desc.CourseDesc(),
		Type:        c.Type.CourseType(),
		Status:      c.Status.CourseStatus(),
		Poster:      c.Poster.URL(),
		Duration:    c.Duration.CourseDuration(),
	}
}
