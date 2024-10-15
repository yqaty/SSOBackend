package handler

import (
	"errors"
	"github.com/UniqueStudio/UniqueSSOBackend/config"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/pkg/open"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type PhoneSMS struct {
	Phone string `json:"phone" binding:"required"`
}

func SendSMSCodeHandler(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "SendValidateCodeSMS")
	defer span.End()

	p := new(PhoneSMS)
	if err := c.ShouldBindJSON(p); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind body failed", zap.Error(err))
		respParamBindingError(c, err)
		return
	}

	smsCode, err := utils.GenerateTmpCode(apmCtx, p.Phone, constants.VerCodeExpireDuration)
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("generate sms code failed", zap.Error(err))
		respDBError(c, err)
		return
	}

	resp, err := core.OpenClient.PushSMS(apmCtx, &open.PushSMSRequest{
		TemplateId: config.Config.OpenPlatform.VerificationCode.Sms.TemplateId,
		Phone:      p.Phone,
		Params:     []string{smsCode},
	})
	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("send verification code failed", zap.Error(err), zap.String("phone", p.Phone))
		respDBError(c, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("open-platform error")
		span.RecordError(err)
		bodyBytes, _ := io.ReadAll(resp.Body)
		zapx.WithContext(apmCtx).With(
			zap.Error(err),
			zap.Int("open status code", resp.StatusCode),
			zap.String("response body", string(bodyBytes))).Error("open resp status code is not 200")
		respDBError(c, err)
		return
	}

	respOK(c, "success")
}
