package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type TrainingCreateCmd struct {
	User      domain.Account
	ProjectId string

	domain.TrainingConfig
}

func (cmd *TrainingCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating training")

	b := cmd.User != nil &&
		cmd.ProjectId != "" &&
		cmd.ProjectName != nil &&
		cmd.ProjectRepoId != "" &&
		cmd.Name != nil &&
		cmd.CodeDir != nil &&
		cmd.BootFile != nil

	if !b {
		return err
	}

	c := &cmd.Compute
	if c.Flavor == nil || c.Type == nil || c.Version == nil {
		return err
	}

	f := func(kv []domain.KeyValue) error {
		for i := range kv {
			if kv[i].Key == nil {
				return err
			}
		}

		return nil
	}

	if f(cmd.Hyperparameters) != nil {
		return err
	}

	if f(cmd.Env) != nil {
		return err
	}

	for i := range cmd.Inputs {
		v := &cmd.Inputs[i]

		if v.Key == nil || v.User == nil || v.Type == nil || v.RepoId == "" {
			return errors.New("invalide input")
		}
	}

	return nil
}

func (cmd *TrainingCreateCmd) toTrainingConfig() *domain.TrainingConfig {
	return &cmd.TrainingConfig
}

type TrainingSummaryDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Error     string `json:"error"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	IsDone    bool   `json:"is_done"`
	Duration  int    `json:"duration"`
}

func (s trainingService) toTrainingSummaryDTO(
	t *domain.TrainingSummary, dto *TrainingSummaryDTO,
) {
	status := t.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	*dto = TrainingSummaryDTO{
		Id:        t.Id,
		Name:      t.Name.TrainingName(),
		Error:     t.Error,
		Status:    status,
		IsDone:    s.isJobDone(t.Status),
		Duration:  t.Duration,
		CreatedAt: utils.ToDate(t.CreatedAt),
	}

	if t.Desc != nil {
		dto.Desc = t.Desc.TrainingDesc()
	}
}

type TrainingDTO struct {
	Id        string `json:"id"`
	ProjectId string `json:"project_id"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	IsDone    bool       `json:"is_done"`
	Error     string     `json:"error"`
	Status    string     `json:"status"`
	Duration  int        `json:"duration"`
	CreatedAt string     `json:"created_at"`
	Compute   ComputeDTO `json:"compute"`
	AimPath   string     `json:"aim_path"`
	EnableAim bool       `json:"enable_aim"`

	LogPreviewURL string `json:"-"`
}

type ComputeDTO struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Flavor  string `json:"flavor"`
}

func (s trainingService) toTrainingDTO(dto *TrainingDTO, ut *domain.UserTraining, link string) {
	t := &ut.TrainingConfig
	detail := &ut.JobDetail
	c := &t.Compute

	status := detail.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	*dto = TrainingDTO{
		Id:        ut.Id,
		ProjectId: ut.ProjectId,

		Name:      t.Name.TrainingName(),
		IsDone:    s.isJobDone(detail.Status),
		Error:     detail.Error,
		Status:    status,
		Duration:  detail.Duration,
		CreatedAt: utils.ToDate(ut.CreatedAt),
		Compute: ComputeDTO{
			Type:    c.Type.ComputeType(),
			Flavor:  c.Flavor.ComputeFlavor(),
			Version: c.Version.ComputeVersion(),
		},
		EnableAim: t.EnableAim,
		AimPath:   ut.JobDetail.AimPath,

		LogPreviewURL: link,
	}

	if t.Desc != nil {
		dto.Desc = t.Desc.TrainingDesc()
	}
}

type TrainingConfigDTO struct {
	Name string `json:"name"`
	Desc string `json:"desc"`

	CodeDir  string `json:"code_dir"`
	BootFile string `json:"boot_file"`

	Hyperparameters []KeyValue    `json:"hyperparameter"`
	Env             []KeyValue    `json:"env"`
	Models          []TrainingRef `json:"models"`
	Datasets        []TrainingRef `json:"datasets"`
	EnableAim       bool          `json:"enable_aim"`
	EnableOutput    bool          `json:"enable_output"`

	Compute Compute `json:"compute"`
}

type Compute struct {
	Type    string `json:"type"`
	Flavor  string `json:"flavor"`
	Version string `json:"version"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TrainingRef struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
	Name  string `json:"name"`
	File  string `json:"File"`
}

func toTrainingConfigDTO(config *domain.TrainingConfig) TrainingConfigDTO {
	return TrainingConfigDTO{
		Name: config.Name.TrainingName(),
		Desc: config.Desc.TrainingDesc(),

		CodeDir:  config.CodeDir.Directory(),
		BootFile: config.BootFile.FilePath(),

		Hyperparameters: toKeyValue(config.Hyperparameters),
		Env: toKeyValue(config.Env),
		Models: toTrainingRef(config.Inputs),

		EnableAim:    config.EnableAim,
		EnableOutput: config.EnableOutput,
	}
}

func toKeyValue(kv []domain.KeyValue) []KeyValue {
	keyValue := make([]KeyValue, len(kv))

	for i := range kv {
		keyValue[i] = KeyValue{
			Key:   kv[i].Key.CustomizedKey(),
			Value: kv[i].Value.CustomizedValue(),
		}
	}

	return keyValue
}

func toTrainingRef(inputs []domain.Input) []TrainingRef {
	trainingRef := make([]TrainingRef, len(inputs))

	for i := range inputs {
		trainingRef[i] = TrainingRef{
			Key:   inputs[i].Key.CustomizedKey(),
			Owner: inputs[i].User.Account(),
			// TODO missing name
			File: inputs[i].File.InputeFilePath(),
		}
	}

	return trainingRef
}
