package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/course/app"
	"github.com/opensourceways/xihe-server/course/domain"
)

func AddRouterForCourseController(
	rg *gin.RouterGroup,
	s app.CourseService,
) {
	ctl := CourseController{
		s: s,
	}

	rg.POST("/v1/course/:id/player", ctl.Apply)
	rg.GET("/v1/course", ctl.List)
}

type CourseController struct {
	baseController

	s app.CourseService
}

// @Summary Apply
// @Description apply the course
// @Tags  Course
// @Param	id	path	string				true	"course id"
// @Param	body body	StudentApplyRequest	true	"body of applying"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/player [post]
func (ctl *CourseController) Apply(ctx *gin.Context) {
	req := StudentApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd, err := req.toCmd(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

		return
	}

	if err := ctl.s.Apply(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary List
// @Description list the course
// @Tags  Course
// @Param	status	query	string	false	"course status, such as over, preparing, in-progress"
// @Param	type	query	string	false	"course type, such as ai, mindspore, foundation"
// @Param	mine	query	string	false	"just list courses of player, if it is set"
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/course [get]
func (ctl *CourseController) List(ctx *gin.Context) {
	var cmd app.CourseListCmd
	var err error

	if str := ctl.getQueryParameter(ctx, "status"); str != "" {
		cmd.Status, err = domain.NewCourseStatus(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	if str := ctl.getQueryParameter(ctx, "type"); str != "" {
		cmd.Type, err = domain.NewCourseType(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}

	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if !visitor && ctl.getQueryParameter(ctx, "mine") != "" {
		cmd.User = pl.DomainAccount()
	}

	if data, err := ctl.s.List(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}
