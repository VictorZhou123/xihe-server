package app

import (
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"github.com/opensourceways/xihe-server/course/domain/user"
)

type CourseService interface {
	// player
	Apply(*PlayerApplyCmd) error
	List(*CourseListCmd) ([]CourseSummaryDTO, error)
}

func NewCourseService(
	userCli user.User,

	courseRepo repository.Course,
	playerRepo repository.Player,
) *courseService {
	return &courseService{
		userCli:    userCli,
		courseRepo: courseRepo,
		playerRepo: playerRepo,
	}
}

type courseService struct {
	userCli user.User

	courseRepo repository.Course
	playerRepo repository.Player
}

// List
func (s *courseService) List(cmd *CourseListCmd) (
	dtos []CourseSummaryDTO, err error,
) {
	return s.listCourses(&repository.CourseListOption{
		Status: cmd.Status,
		Type:   cmd.Type,
	})

}

func (s *courseService) listCourses(opt *repository.CourseListOption) (
	dtos []CourseSummaryDTO, err error,
) {
	v, err := s.courseRepo.FindCourses(opt)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]CourseSummaryDTO, len(v))
	for i := range v {
		n, err := s.playerRepo.PlayerCount(v[i].Id)
		if err != nil {
			return nil, err
		}

		s.toCourseSummaryDTO(&v[i], n, &dtos[i])
	}

	return
}
