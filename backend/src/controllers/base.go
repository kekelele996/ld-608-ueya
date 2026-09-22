package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/middlewares"
	"groundTurn/src/types"
)

func bindID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middlewares.AbortBadRequest(c, constants.MsgNotFound)
		return 0, false
	}
	return uint(id), true
}

// fail lets every controller wrap service-layer AppError itself.
func fail(c *gin.Context, err error) {
	var appErr *types.AppError
	if errors.As(err, &appErr) {
		middlewares.AbortEnvelope(c, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	middlewares.AbortEnvelope(c, http.StatusInternalServerError, constants.CodeInternal, constants.MsgInternal)
}

// MustBindJSON wraps gin binding into the standard validation error envelope.
func MustBindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		middlewares.AbortEnvelope(c, http.StatusBadRequest, constants.CodeValidationFailed,
			"请求参数校验失败："+err.Error())
		return false
	}
	return true
}
