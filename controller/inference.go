package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForInferenceController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	repo repository.Inference,
	project repository.Project,
	sender message.Sender,
) {
	ctl := InferenceController{
		s: app.NewInferenceService(
			p, repo, project, sender, apiConfig.MinSurvivalTimeOfInference,
		),
		project: project,
	}

	rg.POST("/v1/inference/project", ctl.Create)
	rg.GET("/v1/inference/project/:owner/:pid/:instid", ctl.Get)
}

type InferenceController struct {
	baseController

	project repository.Project

	s app.InferenceService
}

// @Summary		Create
// @Description	create inference
// @Tags			Inference
// @Param			body	body	InferenceCreateRequest	true	"body of creating inference"
// @Accept			json
// @Success		201	{object}			app.InferenceDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/inference/project [post]
func (ctl *InferenceController) Create(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	var u domain.Account
	if pl.Account == "" {
		u, _ = domain.NewAccount("unknow")
	} else {
		u = pl.DomainAccount()
	}

	req := InferenceCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd(u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	dto, err := ctl.s.Create(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		log.Errorf("inference failed err:%s", err.Error())

		return
	}

	utils.DoLog("", pl.Account, "create gradio",
		fmt.Sprintf("projectid: %s", cmd.ProjectId), "success")

	ctx.JSON(http.StatusCreated, newResponseData(dto))
}

// @Summary		Get
// @Description	create inference
// @Tags			Inference
// @Param			owner			path	string	true	"project owner"
// @Param			pid				path	string	true	"project id"
// @Param			instid			path	string	true	"inference id"
// @Param			lastcommit		path	string	true	"last commit id"
// @Accept			json
// @Success		201	{object}			app.InferenceDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/inference/project/{owner}/{pid}/{instid} [get]
func (ctl *InferenceController) Get(ctx *gin.Context) {
	_, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, true)
	if !ok {
		return
	}

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{csrftoken},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == csrftoken
		},
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//TODO delete
		log.Errorf("update ws failed, err:%s", err.Error())

		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ws.WriteJSON(newResponseError(err))

		return
	}

	lastCommit, err := ctl.s.GetLastCommitIdByProjectId(&app.GetLastCommitIdByProjectIdCmd{
		User:      owner,
		ProjectId: ctx.Param("pid"),
	})
	if err != nil {
		logrus.Debugf("get last commit id error: %s", err.Error())

		ws.WriteJSON(newResponseError(err))

		return
	}

	// start
	info, err := toInferenceIndex(
		ctx.Param("instid"),
		lastCommit,
		ctx.Param("pid"),
		ctx.Param("owner"),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	// cycle query
	for i := 0; i < apiConfig.InferenceTimeout; i++ {
		dto, err := ctl.s.Get(&info)
		if err != nil {
			ws.WriteJSON(newResponseError(err))

			log.Errorf("inference failed: get status, err:%s", err.Error())

			return
		}

		log.Debugf("info dto:%v", dto)

		if dto.Error != "" || dto.AccessURL != "" {
			ws.WriteJSON(newResponseData(dto))

			log.Debug("inference done")

			return
		}

		time.Sleep(time.Second)
	}

	log.Error("inference timeout")

	ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout"))
}
