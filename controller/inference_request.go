package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type InferenceCreateRequest struct {
	Owner     string `json:"owner"`
	ProjectId string `json:"project_id"`
}

func (req *InferenceCreateRequest) Validate() error {
	if req.ProjectId == "" {
		return errors.New("invalid project id")
	}

	return nil
}

func (req *InferenceCreateRequest) toCmd(user domain.Account) (
	cmd app.InferenceCreateInputCmd, err error,
) {
	if err = req.Validate(); err != nil {
		return
	}

	owner, err := domain.NewAccount(req.Owner)
	if err != nil {
		return
	}

	return app.InferenceCreateInputCmd{
		User:      user,
		Owner:     owner,
		ProjectId: req.ProjectId,
	}, nil
}

func toInferenceIndex(
	instId, lastCommit, pid, owner string,
) (index domain.InferenceIndex, err error) {
	b := instId == "" ||
		lastCommit == "" ||
		pid == ""
	if b {
		err = errors.New("input params error")
	}

	o, err := domain.NewAccount(owner)
	if err != nil {
		return
	}

	return domain.InferenceIndex{
		Project: domain.ResourceIndex{
			Owner: o,
			Id:    pid,
		},
		Id:         instId,
		LastCommit: lastCommit,
	}, nil
}
