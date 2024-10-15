package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UniformResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func respOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, UniformResponse{
		Message: "success",
		Data:    data,
	})
}

func respParamBindingError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, UniformResponse{
		Message: err.Error(),
	})
}

func respExternalAPIError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, UniformResponse{
		Message: err.Error(),
	})
}

func respDBError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, UniformResponse{
		Message: err.Error(),
	})
}

func respForbiddenError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusForbidden, UniformResponse{
		Message: err.Error(),
	})
}

func respUnauthError(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, UniformResponse{
		Message: "permission delay",
	})
}
