package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type InferenceIndex = domain.InferenceIndex
type InferenceDetail = domain.InferenceDetail

type InferenceCreateInputCmd struct {
	User      domain.Account
	Owner     domain.Account
	ProjectId string
}

type inferenceCreateCmd struct {
	ProjectId     string
	ProjectName   domain.ResourceName
	ProjectOwner  domain.Account
	ResourceLevel string

	InferenceDir domain.Directory
	BootFile     domain.FilePath
}

func (cmd *inferenceCreateCmd) Validate() error {
	b := cmd.ProjectId != "" &&
		cmd.ProjectName != nil &&
		cmd.ProjectOwner != nil &&
		cmd.InferenceDir != nil &&
		cmd.BootFile != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *inferenceCreateCmd) toInference(v *domain.Inference, lastCommit, requester string) {
	v.Project.Id = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectName = cmd.ProjectName
	v.ResourceLevel = cmd.ResourceLevel
	v.Project.Owner = cmd.ProjectOwner
	v.Requester = requester
}

func (cmd *inferenceCreateCmd) toRepoDirFile() platform.RepoDirFile {
	return platform.RepoDirFile{
		RepoName: cmd.ProjectName,
		Dir:      cmd.InferenceDir,
		File:     cmd.BootFile,
	}
}

type getLastCommitIdCmd struct {
	User domain.Account
	platform.RepoDirFile
}

type GetLastCommitIdByProjectIdCmd struct {
	User      domain.Account
	ProjectId string
}

type InferenceDTO struct {
	expiry       int64
	Error        string `json:"error"`
	AccessURL    string `json:"access_url"`
	InstanceId   string `json:"inference_id"`
	LastCommitId string `json:"last_commitid"`
}

func (dto *InferenceDTO) hasResult() bool {
	return dto.InstanceId != ""
}

func (dto *InferenceDTO) canReuseCurrent() bool {
	return dto.AccessURL != ""
}
