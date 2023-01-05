package domain

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/utils"
)

const (
	bigmodelVQA         = "vqa"
	bigmodelPanGu       = "pangu"
	bigmodelLuoJia      = "luojia"
	bigmodelWuKong      = "wukong"
	bigmodelCodeGeex    = "codegeex"
	bigmodelGenPicture  = "gen_picture"
	bigmodelDescPicture = "desc_picture"
)

var (
	BigmodelVQA         = BigmodelType(bigmodelVQA)
	BigmodelPanGu       = BigmodelType(bigmodelPanGu)
	BigmodelLuoJia      = BigmodelType(bigmodelLuoJia)
	BigmodelWuKong      = BigmodelType(bigmodelWuKong)
	BigmodelCodeGeex    = BigmodelType(bigmodelCodeGeex)
	BigmodelGenPicture  = BigmodelType(bigmodelGenPicture)
	BigmodelDescPicture = BigmodelType(bigmodelDescPicture)
)

// BigmodelType
type BigmodelType string

// Question
type Question interface {
	Question() string
}

func NewQuestion(v string) (Question, error) {
	// TODO check format

	return question(v), nil
}

type question string

func (s question) Question() string {
	return string(s)
}

// WuKongPictureDesc
type WuKongPictureDesc interface {
	WuKongPictureDesc() string
}

func NewWuKongPictureDesc(v string) (WuKongPictureDesc, error) {
	if v == "" {
		return nil, errors.New("no desc")
	}

	if max := config.WuKongPictureMaxDescLength; utils.StrLen(v) > max {
		return nil, fmt.Errorf(
			"the length of desc should be less than %d", max,
		)
	}

	return wukongPictureDesc(v), nil
}

type wukongPictureDesc string

func (r wukongPictureDesc) WuKongPictureDesc() string {
	return string(r)
}
