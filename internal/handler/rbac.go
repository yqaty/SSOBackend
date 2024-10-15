package handler

import (
	"errors"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/model"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

type ObjectInfo struct {
	ID       uint   `json:"id"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type CheckPermissionByRoleRequest struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
}

type CheckPermissionResponse struct {
	Ok bool `json:"ok"`
}

func GetUserInfoByUID(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "GetUserInfoByUid")
	defer span.End()

	requestUID, ok := c.GetQuery("uid")
	if !ok || requestUID == "" {
		err := errors.New("query: uid field not found")
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind query param uid failed", zap.String("uid", requestUID))
		respParamBindingError(c, err)
		return
	}

	user, err := model.GetUserByUID(apmCtx, requestUID)
	if err != nil {
		span.RecordError(err)
		zapx.With(zap.Error(err), zap.String("uid", requestUID)).Error("GetUserByUID error")
		respDBError(c, err)
		return
	}

	roles, err := model.GetRolesByUID(apmCtx, requestUID)
	if err != nil {
		span.RecordError(err)
		respDBError(c, err)
		return
	}

	respOK(c, UserDetail{
		UID:       user.UID,
		Phone:     user.Phone,
		Email:     user.Email,
		Roles:     roles,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		Gender:    user.Gender,
		Groups:    user.Groups,
	})
}

func GetUserRolesHandler(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "GetUserRoleHandler")
	defer span.End()

	uid := mustGetUIDFromCtx(apmCtx)

	requestUID, ok := c.GetQuery("uid")
	if !ok || requestUID == "" {
		err := errors.New("query: uid field not found")
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind query param uid failed", zap.String("uid", uid))
		respParamBindingError(c, err)
		return
	}

	roles, err := model.GetRolesByUID(apmCtx, requestUID)
	if err != nil {
		span.RecordError(err)
		respDBError(c, err)
		return
	}

	respOK(c, gin.H{requestUID: roles})
}

func GetAllRolesHandler(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "GetAllRoleHandler")
	defer span.End()

	roles, err := model.GetAllRoles(apmCtx)
	if err != nil {
		span.RecordError(err)
		respDBError(c, err)
		return
	}

	respOK(c, gin.H{"roles": roles})
}
