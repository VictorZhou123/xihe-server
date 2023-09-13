package messageadapter

import (
	"strconv"
	"strings"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	basemsg "github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/utils"
)

func NewPublisher(cfg *Config) *publisher {
	return &publisher{*cfg}
}

type publisher struct {
	cfg Config
}

func (impl *publisher) SendWuKongInferenceStart(v *domain.WuKongInferenceStartEvent) error {
	cfg := &impl.cfg.InferenceStart

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"status":    "waiting",
			"task_type": v.EsStyle,
			"style":     v.Style,
			"desc":      v.Desc.WuKongPictureDesc(),
		},
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

func (impl *publisher) SendWuKongInferenceError(v *domain.WuKongInferenceErrorEvent) error {
	cfg := &impl.cfg.InferenceError

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"task_id": strconv.Itoa(int(v.TaskId)),
			"status":  "error",
			"error":   v.ErrMsg,
		},
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

func (impl *publisher) SendWuKongAsyncTaskStart(v *domain.WuKongAsyncTaskStartEvent) error {
	cfg := &impl.cfg.InferenceAsyncStart

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"status":  "running",
			"task_id": strconv.Itoa(int(v.TaskId)),
		},
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

func (impl *publisher) SendWuKongAsyncInferenceFinish(
	v *domain.WuKongAsyncInferenceFinishEvent,
) error {
	cfg := &impl.cfg.InferenceAsyncFinish

	var ls string
	for k := range v.Links { // TODO: Move it into domain.service
		ls += v.Links[k] + ","
	}

	msg := common.MsgNormal{
		User: v.Account.Account(),
		Details: map[string]string{
			"task_id": strconv.Itoa(int(v.TaskId)),
			"status":  "finished",
			"links":   strings.TrimRight(ls, ","),
		},
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

func (impl *publisher) SendBigModelAccessLog(v *domain.BigModelAccessLogEvent) error {
	cfg := &impl.cfg.BigModelAccessLog

	msg := basemsg.MsgOperateLog{
		When: utils.Now(),
		User: v.Account.Account(),
		Type: "bigmodel",
		Info: map[string]string{
			"bigmodel": string(v.BigModelType),
		},
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

func (impl *publisher) SendWuKongPicturePublic(v *domain.WuKongPicturePublicEvent) error {
	cfg := &impl.cfg.PicturePublic

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		CreatedAt: utils.Now(),
	}

	return common.Publish(cfg.Topic, &msg, nil)
}

// Config
type Config struct {
	// wukong
	InferenceStart       common.TopicConfig `json:"inference_start"`
	InferenceError       common.TopicConfig `json:"inference_error"`
	InferenceAsyncStart  common.TopicConfig `json:"inference_async_start"`
	InferenceAsyncFinish common.TopicConfig `json:"inference_async_finish"`
	PicturePublic        common.TopicConfig `json:"picture_public"`

	// common
	BigModelAccessLog common.TopicConfig `json:"bigmodel_access_log"`
}
