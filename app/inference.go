package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type InferenceService interface {
	Create(*InferenceCreateInputCmd) (InferenceDTO, error)
	Get(info *InferenceIndex) (InferenceDTO, error)
	GetLastCommitIdByProjectId(*GetLastCommitIdByProjectIdCmd) (string, error)
}

func NewInferenceService(
	p platform.RepoFile,
	repo repository.Inference,
	projectRepo repository.Project,
	sender message.Sender,
	minSurvivalTime int,
) InferenceService {
	return inferenceService{
		p:               p,
		repo:            repo,
		projectRepo:     projectRepo,
		sender:          sender,
		minSurvivalTime: int64(minSurvivalTime),
	}
}

type inferenceService struct {
	p               platform.RepoFile
	repo            repository.Inference
	projectRepo     repository.Project
	sender          message.Sender
	minSurvivalTime int64
}

func (s inferenceService) Create(cmd *InferenceCreateInputCmd) (dto InferenceDTO, err error) {
	// get project summary
	v, err := s.projectRepo.GetSummary(cmd.Owner, cmd.ProjectId)
	if err != nil {
		return
	}

	// is private
	if v.IsPrivate() {
		err = errors.New("project is private")

		return
	}

	// get resource level
	var level string
	if level, err = s.getResourceLevel(v.Owner, v.Id); err != nil {
		return
	}

	// create
	inferenceDir, _ := domain.NewDirectory(appConfig.InferenceDir)
	inferenceBootFile, _ := domain.NewFilePath(appConfig.InferenceBootFile)
	inferCreateCmd := inferenceCreateCmd{
		ProjectId:     v.Id,
		ProjectName:   v.Name,
		ProjectOwner:  v.Owner,
		ResourceLevel: level,
		InferenceDir:  inferenceDir,
		BootFile:      inferenceBootFile,
	}

	dto, err = s.create(cmd.User.Account(), v.Owner, &inferCreateCmd)
	if err != nil {
		return
	}

	return
}

func (s inferenceService) create(user string, owner domain.Account, cmd *inferenceCreateCmd) (
	dto InferenceDTO, err error,
) {
	sha, err := s.getLastCommitId(&getLastCommitIdCmd{
		User:        owner,
		RepoDirFile: cmd.toRepoDirFile(),
	})
	if err != nil {
		return
	}

	instance := new(domain.Inference)
	cmd.toInference(instance, sha, user)

	dto, version, err := s.check(instance)
	if err != nil {
		return
	}

	if dto.hasResult() {
		if dto.canReuseCurrent() {
			instance.Id = dto.InstanceId
			logrus.Debugf("will reuse the inference instance(%s)", dto.InstanceId)

			err1 := s.sender.ExtendInferenceSurvivalTime(&message.InferenceExtendInfo{
				InferenceInfo: instance.InferenceInfo,
				Expiry:        dto.expiry,
			})
			if err1 != nil {
				logrus.Errorf(
					"extend instance(%s) failed, err:%s",
					dto.InstanceId, err1.Error(),
				)
			}
		}

		return
	}

	if dto.InstanceId, err = s.repo.Save(instance, version); err == nil {
		instance.Id = dto.InstanceId
		err = s.sender.CreateInference(&instance.InferenceInfo)

		return
	}

	if repository.IsErrorDuplicateCreating(err) {
		dto, _, err = s.check(instance)
	}

	return
}

func (s inferenceService) Get(index *InferenceIndex) (dto InferenceDTO, err error) {
	v, err := s.repo.FindInstance(index)
	if err != nil {
		return
	}

	dto.Error = v.Error
	dto.AccessURL = v.AccessURL
	dto.InstanceId = v.Id

	return
}

func (s inferenceService) check(instance *domain.Inference) (
	dto InferenceDTO, version int, err error,
) {
	v, version, err := s.repo.FindInstances(&instance.Project, instance.LastCommit)
	if err != nil || len(v) == 0 {
		return
	}

	var target *repository.InferenceSummary

	for i := range v {
		item := &v[i]

		if item.Error != "" {
			dto.Error = item.Error
			dto.InstanceId = item.Id

			return
		}

		if target == nil || item.Expiry > target.Expiry {
			target = item
		}
	}

	if target == nil {
		return
	}

	e, n := target.Expiry, utils.Now()
	if n < e && n+s.minSurvivalTime <= e {
		dto.expiry = target.Expiry
		dto.AccessURL = target.AccessURL
		dto.InstanceId = target.Id
	}

	return
}

func (s inferenceService) getResourceLevel(owner domain.Account, pid string) (level string, err error) {
	resources, err := s.projectRepo.FindUserProjects(
		[]repository.UserResourceListOption{
			{
				Owner: owner,
				Ids: []string{
					pid,
				},
			},
		},
	)

	if err != nil || len(resources) < 1 {
		return
	}

	if resources[0].Level != nil {
		level = resources[0].Level.ResourceLevel()
	}

	return
}

func (s inferenceService) getLastCommitId(cmd *getLastCommitIdCmd) (string, error) {
	u := platform.UserInfo{
		User: cmd.User,
	}

	sha, b, err := s.p.GetDirFileInfo(&u, &cmd.RepoDirFile)
	if err != nil {
		return "", err
	}

	if !b {
		err = ErrorUnavailableRepoFile{
			errors.New("no boot file"),
		}
	}

	return sha, nil
}

func (s inferenceService) GetLastCommitIdByProjectId(cmd *GetLastCommitIdByProjectIdCmd) (string, error) {
	// get project name
	p, err := s.projectRepo.Get(cmd.User, cmd.ProjectId)
	if err != nil {
		return "", err
	}

	inferenceDir, _ := domain.NewDirectory(appConfig.InferenceDir)
	inferenceBootFile, _ := domain.NewFilePath(appConfig.InferenceBootFile)

	return s.getLastCommitId(&getLastCommitIdCmd{
		User: cmd.User,
		RepoDirFile: platform.RepoDirFile{
			RepoName: p.Name,
			Dir:      inferenceDir,
			File:     inferenceBootFile,
		},
	})
}

type InferenceInternalService interface {
	UpdateDetail(*InferenceIndex, *InferenceDetail) error
}

func NewInferenceInternalService(repo repository.Inference) InferenceInternalService {
	return inferenceInternalService{
		repo: repo,
	}
}

type inferenceInternalService struct {
	repo repository.Inference
}

func (s inferenceInternalService) UpdateDetail(index *InferenceIndex, detail *InferenceDetail) error {
	return s.repo.UpdateDetail(index, detail)
}

type InferenceMessageService interface {
	CreateInferenceInstance(*domain.InferenceInfo) error
	ExtendSurvivalTime(*message.InferenceExtendInfo) error
}

func NewInferenceMessageService(
	repo repository.Inference,
	user userrepo.User,
	manager inference.Inference,
) InferenceMessageService {
	return inferenceMessageService{
		repo:    repo,
		user:    user,
		manager: manager,
	}
}

type inferenceMessageService struct {
	repo    repository.Inference
	user    userrepo.User
	manager inference.Inference
}

func (s inferenceMessageService) CreateInferenceInstance(info *domain.InferenceInfo) error {
	v, err := s.user.GetByAccount(info.Project.Owner)
	if err != nil {
		return err
	}

	survivaltime, err := s.manager.Create(&inference.InferenceInfo{
		InferenceInfo: info,
		UserToken:     v.PlatformToken,
	})
	if err != nil {
		return err
	}

	return s.repo.UpdateDetail(
		&info.InferenceIndex,
		&domain.InferenceDetail{Expiry: utils.Now() + int64(survivaltime)},
	)
}

func (s inferenceMessageService) ExtendSurvivalTime(info *message.InferenceExtendInfo) error {
	expiry, n := info.Expiry, utils.Now()
	if expiry < n {
		logrus.Errorf(
			"extend survival time for inference instance(%s) failed, it is timeout.",
			info.Id,
		)

		return nil
	}

	n += int64(s.manager.GetSurvivalTime(&info.InferenceInfo))

	v := int(n - expiry)
	if v < 10 {
		logrus.Debugf(
			"no need to extend survival time for inference instance(%s) in a small range",
			info.Id,
		)

		return nil
	}

	if err := s.manager.ExtendSurvivalTime(&info.InferenceIndex, v); err != nil {
		return err
	}

	return s.repo.UpdateDetail(&info.InferenceIndex, &domain.InferenceDetail{Expiry: n})
}
