package handler

import (
	"errors"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/model"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

type ResetPasswordRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func ResetPassword(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "UpdateUserInfoHandler")
	defer span.End()
	zlog := zapx.WithContext(apmCtx)

	var req ResetPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		zlog.With(zap.Error(err)).Error("bind json error")
		respParamBindingError(c, err)
		return
	}

	tempCode, err := utils.GetTmpCodeByID(apmCtx, req.Phone)
	if err != nil {
		zlog.With(zap.Error(err)).Error("GetTmpCodeByID error")
		respDBError(c, err)
		return
	}

	if tempCode != req.Code {
		err = errors.New("validate code wrong")
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("", zap.Error(err))
		respForbiddenError(c, err)
		return
	}

	err = model.ResetPassword(apmCtx, req.Phone, req.Password)
	if err != nil {
		zlog.With(zap.Error(err)).Error("ResetPassword error")
		respDBError(c, err)
		return
	}

	respOK(c, nil)
}
