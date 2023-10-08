package bigmodels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/sirupsen/logrus"
	"github.com/victorzhou123/sse"
)

type glm2Request struct {
	Inputs            string      `json:"inputs"`
	History           [][2]string `json:"history"`
	Sampling          bool        `json:"sampling"`
	TopK              int         `json:"top_k"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	RepetitionPenalty float64     `json:"repetition_penalty"`
}

type glm2Response struct {
	Reply        string `json:"reply"`
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	StreamStatus string `json:"stream_status"`
}

type glm2Info struct {
	endpoints chan string
}

func newGLM2Info(cfg *Config) (info glm2Info, err error) {
	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.GLM2)

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	return
}

func (s *service) GLM2(ch chan string, input *domain.GLM2Input) (err error) {
	// input audit
	if err = s.check.check(input.Text.GLM2Text()); err != nil {
		return
	}

	// call bigmodel glm2
	f := func(ec chan string, e string) (err error) {
		err = s.genGLM2(ec, ch, e, input)

		return
	}

	if err = s.doIfFreeNoEndpointReturn(s.glm2Info.endpoints, f); err != nil {
		return
	}

	return
}

// func (s *service) genGLM2(ec, ch chan string, endpoint string, input *domain.GLM2Input) (
// 	err error,
// ) {
// 	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
// 	if err != nil {
// 		return
// 	}

// 	opt := toGLM2Req(input)
// 	body, err := libutils.JsonMarshal(&opt)
// 	if err != nil {
// 		return
// 	}

// 	req, err := http.NewRequest(
// 		http.MethodPost, endpoint, bytes.NewBuffer(body),
// 	)
// 	if err != nil {
// 		return
// 	}

// 	req.Header.Set("X-Auth-Token", t)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Connection", "keep-alive")
// 	req.Header.Set("Accept", "*/*")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return
// 	}

// 	reader := bufio.NewReader(resp.Body)

// 	var r glm2Response
// 	go func() {
// 		defer func() { ec <- endpoint }()
// 		defer resp.Body.Close()

// 		for {
// 			line, err := reader.ReadString('\n')
// 			if err != nil {
// 				ch <- "done"

// 				return
// 			}

// 			fmt.Printf("line: %v\n", line)

// 			data := strings.Replace(string(line), "data: ", "", 1)
// 			data = strings.TrimRight(data, "\x00")

// 			if err = json.Unmarshal([]byte(data), &r); err != nil {
// 				continue
// 			}

// 			if r.StreamStatus == "DONE" {
// 				ch <- "done"

// 				return
// 			}

// 			// response audit
// 			if r.Reply != "" {
// 				if err = s.check.check(r.Reply); err != nil {
// 					logrus.Debugf("content audit not pass: %s", err.Error())

// 					ch <- "done"

// 					return
// 				}
// 			}

// 			ch <- r.Reply
// 		}
// 	}()

// 	return
// }

func (s *service) genGLM2(ec, ch chan string, endpoint string, input *domain.GLM2Input) (
	err error,
) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return
	}

	opt := toGLM2Req(input)
	body, err := libutils.JsonMarshal(&opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	cli := sse.NewClient(endpoint, func(c *sse.Client) {
		c.Headers = map[string]string{
			"X-Auth-Token": t,
			"Content-Type": "application/json",
			"Connection":   "keep-alive",
		}

		c.Request = req
		c.Connection = http.DefaultClient
		c.OnDisconnect(func(c *sse.Client) {
			ec <- endpoint
		})
	})

	// if err = ; err != nil {
	// 	fmt.Printf("err2: %v\n", err)
	// 	ch <- "done"
	// }
	go cli.Subscribe("", func(msg *sse.Event) {
		fmt.Printf("string(msg.Data): %v\n", string(msg.Data))
		ch <- s.toReply(string(msg.Data))
	})

	return
}

func (s *service) toReply(v string) string {
	data := strings.TrimPrefix(v, "data: ")
	data = strings.TrimRight(data, "\x00")

	var r glm2Response
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return ""
	}

	if r.StreamStatus == "DONE" {
		return "done"
	}

	// response audit
	if r.Reply != "" {
		if err := s.check.check(r.Reply); err != nil {
			logrus.Debugf("content audit not pass: %s", err.Error())

			return "done"
		}
	}

	return ""
}

func toGLM2Req(input *domain.GLM2Input) glm2Request {
	history := make([][2]string, len(input.History))

	for i := range input.History {
		history[i] = input.History[i].History()
	}

	return glm2Request{
		Inputs:            input.Text.GLM2Text(),
		History:           history,
		Sampling:          input.Sampling,
		TopK:              input.TopK.TopK(),
		TopP:              input.TopP.TopP(),
		Temperature:       input.Temperature.Temperature(),
		RepetitionPenalty: input.RepetitionPenalty.RepetitionPenalty(),
	}
}
