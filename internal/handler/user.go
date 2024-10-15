package handler

import (
	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/model"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/tracer"
	"github.com/gin-gonic/gin"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

type UserDetail struct {
	UID         string           `json:"uid"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	Password    string           `json:"password,omitempty"`
	Roles       []string         `json:"roles"`
	Name        string           `json:"name"`
	AvatarURL   string           `json:"avatar_url"`
	Gender      constants.Gender `json:"gender"`
	JoinTime    string           `json:"join_time"`
	Groups      []string         `json:"groups"`
	LarkUnionID string           `json:"lark_union_id"`
}

func GetUserInfoHandler(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "GetUserInfoHandler")
	defer span.End()

	uid := mustGetUIDFromCtx(apmCtx)
	user, err := model.GetUserByUID(apmCtx, uid)
	if err != nil {
		span.RecordError(err)
		respDBError(c, err)
		return
	}

	roles, err := model.GetRolesByUID(apmCtx, uid)
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
		JoinTime:  user.JoinTime,
		Groups:    user.Groups,
	})
}

type UpdateUser struct {
	//	Phone     string           `json:"phone"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	Gender    constants.Gender `json:"gender"`
	AvatarURL string           `json:"avatar_url"`
}

func UpdateUserInfoHandler(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "UpdateUserInfoHandler")
	defer span.End()

	uid := mustGetUIDFromCtx(apmCtx)

	updateUser := UpdateUser{}
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		span.RecordError(err)
		zapx.WithContext(apmCtx).Error("bind update user failed", zap.Error(err), zap.String("uid", uid))
	}

	if err := model.UpdateUserInfo(apmCtx, &model.User{
		UID:  uid,
		Name: updateUser.Name,
		//	Phone:     updateUser.Phone,
		Email:     updateUser.Email,
		Password:  updateUser.Password,
		Gender:    updateUser.Gender,
		AvatarURL: updateUser.AvatarURL,
	}); err != nil {
		span.RecordError(err)
		respDBError(c, err)
		return
	}

	respOK(c, "success")
}
