package handler

import (
	"context"
	"errors"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"

	"github.com/UniqueStudio/UniqueSSOBackend/config"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/model"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

type LoginUser struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	ValidateCode string `json:"validate_code"`
}

func LoginHandler(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "Login")
	defer span.End()

	lu := new(LoginUser)
	if err := c.ShouldBindJSON(lu); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind param failed", zap.Error(err))
		respParamBindingError(c, err)
		return
	}

	var err error
	var uid string
	switch {
	case lu.Phone != "" && lu.ValidateCode != "":
		uid, err = smsLogin(apmCtx, lu)
	case lu.Email != "" && lu.ValidateCode != "":
		uid, err = emailCodeLogin(apmCtx, lu)
	case lu.Phone != "" && lu.Password != "":
		uid, err = phoneLogin(apmCtx, lu)
	case lu.Email != "" && lu.Password != "":
		uid, err = emailLogin(apmCtx, lu)
	default:
		err = errors.New("unsupported login type")
	}

	if err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("login failed", zap.Error(err))
		respParamBindingError(c, err)
		return
	}

	// login success. add session
	sess := sessions.Default(c)
	// use uid as unique key
	// this is only used for checking whether user has logged in
	sess.Set(constants.SessionNameUID, uid)
	sess.Options(config.SessionOptions)
	//	c.SetCookie("uid", uid, constants.SessionMaxAgeSeconds, "/", "hustunique.com", false, true)
	sess.Set("Path", "/") // set cookie path, because use session.Option set cookie path failed....
	sess.Save()
	respOK(c, "success")
}

func emailLogin(ctx context.Context, lu *LoginUser) (string, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "EmailLogin")
	defer span.End()
	user, err := model.GetUserByEmail(apmCtx, lu.Email)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	if err := utils.ValidatePassword(lu.Password, user.Password); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("validate password failed", zap.Error(err), zap.String("UID", user.UID))
		return "", err
	}

	return user.UID, nil
}

func emailCodeLogin(ctx context.Context, lu *LoginUser) (string, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "EmailCodeLogin")
	defer span.End()

	user, err := model.GetUserByEmail(apmCtx, lu.Email)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	emailCode, err := utils.GetTmpCodeByID(apmCtx, lu.Email)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	if emailCode != lu.ValidateCode {
		err = errors.New("email code is wrong")
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("email code is not match generated one", zap.String("UID", user.UID))
		return "", err
	}

	return user.UID, nil
}

func phoneLogin(ctx context.Context, lu *LoginUser) (string, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "PhoneLogin")
	defer span.End()
	user, err := model.GetUserByPhone(apmCtx, lu.Phone)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	if err := utils.ValidatePassword(lu.Password, user.Password); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("validate password failed", zap.Error(err), zap.String("UID", user.UID))
		return "", err
	}

	return user.UID, nil
}

func smsLogin(ctx context.Context, lu *LoginUser) (string, error) {
	apmCtx, span := tracer.Tracer.Start(ctx, "SMSLogin")
	defer span.End()

	user, err := model.GetUserByPhone(apmCtx, lu.Phone)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	smsCode, err := utils.GetTmpCodeByID(apmCtx, lu.Phone)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	if smsCode != lu.ValidateCode {
		err = errors.New("sms code is wrong")
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("sms code is not match generated one", zap.String("UID", user.UID))
		return "", err
	}

	return user.UID, nil
}
