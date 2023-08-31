package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opensourceways/xihe-server/utils"
)

const (
	bigmodelVQA           = "vqa"
	bigmodelPanGu         = "pangu"
	bigmodelLuoJia        = "luojia"
	bigmodelWuKong        = "wukong"
	bigmodelWuKong4Img    = "wukong_4img"
	bigmodelWuKongHF      = "wukong_hf"
	bigmodelWuKongUser    = "wukong_user"
	bigmodelCodeGeex      = "codegeex"
	bigmodelGenPicture    = "gen_picture"
	bigmodelDescPicture   = "desc_picture"
	bigmodelDescPictureHF = "desc_picture_hf"
	bigmodelAIDetector    = "ai_detector"

	langZH = "zh"
	langEN = "en"

	modelNameWukong = "wukong"
)

var (
	BigmodelVQA           = BigmodelType(bigmodelVQA)
	BigmodelPanGu         = BigmodelType(bigmodelPanGu)
	BigmodelLuoJia        = BigmodelType(bigmodelLuoJia)
	BigmodelWuKong        = BigmodelType(bigmodelWuKong)
	BigmodelWuKong4Img    = BigmodelType(bigmodelWuKong4Img)
	BigmodelWuKongHF      = BigmodelType(bigmodelWuKongHF)
	BigmodelWuKongUser    = BigmodelType(bigmodelWuKongUser)
	BigmodelCodeGeex      = BigmodelType(bigmodelCodeGeex)
	BigmodelGenPicture    = BigmodelType(bigmodelGenPicture)
	BigmodelDescPicture   = BigmodelType(bigmodelDescPicture)
	BigmodelDescPictureHF = BigmodelType(bigmodelDescPictureHF)
	BigmodelAIDetector    = BigmodelType(bigmodelAIDetector)

	wukongPictureLevelMap = map[string]int{
		"official": 2,
		"good":     1,
		"normal":   0,
	}

	resourceLevelMap = map[string]int{
		"official": 2,
		"good":     1,
	}
)

// BigmodelType
type BigmodelType string

// Question
type Question interface {
	Question() string
}

func NewQuestion(v string) (Question, error) {
	// TODO check format
	if v == "" || utils.StrLen(v) > 30 {
		return nil, errors.New("invalid question")
	}

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

	v = utils.XSSFilter(v)

	if max := 55; utils.StrLen(v) > max { // TODO config
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

// wukong level
type WuKongPictureLevel interface {
	WuKongPictureLevel() string
	Int() int
	IsOfficial() bool
}

func NewWuKongPictureLevel(v string) WuKongPictureLevel {
	for k, n := range wukongPictureLevelMap {
		if k == v {
			return wukongPictureLevel{
				level: n,
				desc:  k,
			}
		}
	}

	return nil

}

func NewWuKongPictureLevelByNum(v int) WuKongPictureLevel {
	for k, n := range wukongPictureLevelMap {
		if n == v {
			return wukongPictureLevel{
				level: n,
				desc:  k,
			}
		}
	}

	return nil
}

type wukongPictureLevel struct {
	level int
	desc  string
}

func (r wukongPictureLevel) WuKongPictureLevel() string {
	return r.desc
}

func (r wukongPictureLevel) Int() int {
	return r.level
}

func (r wukongPictureLevel) IsOfficial() bool {
	return r.WuKongPictureLevel() == "official"
}

// obspath
type OBSPath interface {
	OBSPath() string
	IsTempPath() bool
}

func NewOBSPath(v string) (OBSPath, error) {
	if !utils.IsPath(v) {
		return nil, errors.New("invalid obspath")
	}

	return obspath(v), nil
}

type obspath string

func (r obspath) OBSPath() string {
	return string(r)
}

func (r obspath) IsTempPath() bool {
	return strings.Contains(r.OBSPath(), "generate/")
}

type AIDetectorText interface {
	AIDetectorText() string
}

func NewAIDetectorText(v string) (AIDetectorText, error) {
	if v == "" {
		return nil, errors.New("invalid AI detector text")
	}

	v = utils.XSSFilter(v)

	return aidetectortext(v), nil
}

type aidetectortext string

func (r aidetectortext) AIDetectorText() string {
	return string(r)
}

type Lang interface {
	Lang() string
	IsZH() bool
	IsEN() bool
}

func NewLang(v string) (Lang, error) {
	b := v == langZH ||
		v == langEN

	if !b {
		return nil, errors.New("language invalid")
	}

	return lang(v), nil
}

type lang string

func (r lang) Lang() string {
	return string(r)
}

func (r lang) IsZH() bool {
	return r.Lang() == langZH
}

func (r lang) IsEN() bool {
	return r.Lang() == langEN
}

// taichu
type Desc interface {
	Desc() string
}

func NewDesc(v string) (Desc, error) {
	v = utils.XSSFilter(v)

	b := v == "" ||
		utils.StrLen(v) > 30

	if b {
		return nil, errors.New("invalid desc")
	}

	return desc(v), nil
}

type desc string

func (r desc) Desc() string {
	return string(r)
}

// Model Name
type ModelName interface {
	ModelName() string
}

func NewModelName(v string) (ModelName, error) {
	b := v == modelNameWukong
	if !b {
		return nil, errors.New("invalid model name")
	}
	return modelName(v), nil
}

type modelName string

func (m modelName) ModelName() string {
	return string(m)
}

// glm text
type GLMText interface {
	GLMText() string
}

type gLMText string

func NewGLMText(v string) (GLMText, error) {
	if v == "" {
		return nil, errors.New("no glm text")
	}

	if max:= 1000; utils.StrLen(v) > max { // TODO: to config
		return nil, errors.New("invalid glm text")
	}

	return gLMText(v), nil
}

func (g gLMText) GLMText() string {
	return string(g)
}

// top k
type TopK interface {
	TopK() int
}

type topk int

func NewTopK(v int) (TopK, error) {
	if min,max:=0,10; v < min || v > max {
		return nil, errors.New("invalid topk")
	}

	return topk(v), nil
}

func (t topk) TopK() int {
	return int(t)
}

// top p
type TopP interface {
	TopP() float64
}

type topp float64

func NewTopP(v float64) (TopP, error) {
	if min, max := 0.00, 1.00; v < min || v > max {
		return nil, errors.New("invalid topp")
	}

	return topp(v), nil
}

func (t topp) TopP() float64 {
	return float64(t)
}

// temperature
type Temperature interface {
	Temperature() float64
}

type temperature float64

func NewTemperature(v float64) (Temperature, error) {
	if min:= 0.00; v < min {
		return nil, errors.New("invalid temperature")
	}

	return temperature(v), nil
}

func (t temperature) Temperature() float64 {
	return float64(t)
}

// repetition penalty
type RepetitionPenalty interface {
	RepetitionPenalty() float64
}

type repetitionPenalty float64

func NewRepetitionPenalty(v float64) (RepetitionPenalty, error) {
	if min:= 0.00; v < min {
		return nil, errors.New("invalid repetitionPenalty")
	}

	return repetitionPenalty(v), nil
}

func (t repetitionPenalty) RepetitionPenalty() float64 {
	return float64(t)
} 