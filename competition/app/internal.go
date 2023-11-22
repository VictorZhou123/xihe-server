package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
)

type CompetitionSubmissionUpdateCmd = domain.SubmissionUpdatingInfo

// Internal Service
type CompetitionInternalService interface {
	UpdateSubmission(*CompetitionSubmissionUpdateCmd) error
	GetTeamMembers(string) (CompetitionTeamDTO, error)
}

func NewCompetitionInternalService(repo repository.Work, playerRepo repository.Player) CompetitionInternalService {
	return competitionInternalService{
		repo:       repo,
		playerRepo: playerRepo,
	}
}

type competitionInternalService struct {
	repo       repository.Work
	playerRepo repository.Player
}

func (s competitionInternalService) UpdateSubmission(cmd *CompetitionSubmissionUpdateCmd) error {
	w, _, err := s.repo.FindWork(cmd.Index, cmd.Phase)
	if err != nil {
		return err
	}

	submission := w.UpdateSubmission(cmd)
	if submission == nil {
		return errors.New("no corresponding submission")
	}

	v := domain.PhaseSubmission{
		Phase:      cmd.Phase,
		Submission: *submission,
	}

	return s.repo.SaveSubmission(&w, &v)
}

func (s competitionInternalService) GetTeamMembers(repo string) (dto CompetitionTeamDTO, err error) {
	// get work
	work, err := s.repo.FindWorkByRepo(repo)
	if err != nil {
		return
	}

	// get player
	player, err := s.playerRepo.FindPlayerById(work.PlayerId)
	if err != nil {
		return
	}

	s.toCompetitionTeamDTO(&player, &dto)

	return
}
