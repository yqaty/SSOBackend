package model

import (
	"context"
	"time"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRole struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	RoleName string `gorm:"column:role_name;primaryKey"`
	UID      string `gorm:"column:uid;primaryKey"`
}

func (*UserRole) TableName() string { return "user_role" }

func GrantUserRoles(ctx context.Context, uid string, roles []string) error {
	db := core.DB.WithContext(ctx)
	if err := grantUserRoles(db, uid, roles); err != nil {
		zapx.Error("grant user roles failed", zap.Error(err), zap.String("uid", uid), zap.Strings("roles", roles))
		return err
	}
	return nil
}

func AddRoles(ctx context.Context, roles []string) error {
	db := core.DB.WithContext(ctx)
	if err := grantUserRoles(db, "", roles); err != nil {
		zapx.Error("add roles failed", zap.Error(err), zap.Strings("roles", roles))
		return err
	}
	return nil
}

func grantUserRoles(tx *gorm.DB, uid string, roles []string) error {
	ur := make([]UserRole, len(roles))
	for i := range roles {
		ur[i].RoleName = roles[i]
		ur[i].UID = uid
	}
	return tx.Create(ur).Error
}

func RemoveUserRoles(ctx context.Context, uid string, roles []string) error {
	db := core.DB.WithContext(ctx)
	if err := db.Where("uid = ? AND role_name IN ?", uid, roles).Delete(&UserRole{}).Error; err != nil {
		zapx.WithContext(ctx).Error("delete user role failed", zap.Error(err), zap.String("uid", uid), zap.Strings("roles", roles))
		return err
	}
	return nil
}

func GetAllRoles(ctx context.Context) ([]string, error) {
	db := core.DB.WithContext(ctx)
	roles := []string{}
	if err := db.Model(&UserRole{}).Select("role_name").Distinct("role_name").Find(&roles).Error; err != nil {
		zapx.WithContext(ctx).Error("get all roles failed", zap.Error(err))
		return nil, err
	}
	return roles, nil
}

func GetRolesByUID(ctx context.Context, uid string) ([]string, error) {
	db := core.DB.WithContext(ctx)
	roles, err := getRolesByUID(db, uid)
	if err != nil {
		zapx.WithContext(ctx).Error("get role by uid failed", zap.Error(err), zap.String("uid", uid))
		return nil, err
	}
	return roles, nil
}

func getRolesByUID(tx *gorm.DB, uid string) ([]string, error) {
	ur := []UserRole{}
	if err := tx.Where(&UserRole{UID: uid}).Find(&ur).Error; err != nil {
		return nil, err
	}
	roles := make([]string, len(ur))
	for i := range ur {
		roles[i] = ur[i].RoleName
	}
	return roles, nil
}
