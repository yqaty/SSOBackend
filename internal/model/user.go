package model

import (
	"context"
	"errors"
	"time"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/utils"
	"github.com/lib/pq"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrNilUID      = errors.New("uid is nil")
	ErrNilPassword = errors.New("password can not be nil")
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	UID      string `gorm:"column:uid;primaryKey"`
	Phone    string `gorm:"column:phone"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`

	Name string `gorm:"column:name"`
	//JoinTime  time.Time        `gorm:"column:join_time"`
	JoinTime  string           `gorm:"column:join_time"`
	AvatarURL string           `gorm:"column:avatar_url"`
	Gender    constants.Gender `gorm:"column:gender"`

	Groups pq.StringArray `gorm:"column:groups;type:text[]"`
}

func (*User) TableName() string { return "user" }

// AddUser - insert user with auto encrypt password and grant roles
func AddUser(ctx context.Context, user *User, roles ...string) (string, error) {
	if user.Password == "" {
		return "", ErrNilPassword
	}

	encryptedPassword, err := utils.EncryptPassword(user.Password)
	if err != nil {
		zapx.WithContext(ctx).Error("encrypt password failed", zap.Error(err))
		return "", err
	}
	user.Password = encryptedPassword

	db := core.DB.WithContext(ctx)
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			zapx.WithContext(ctx).Error("create user failed", zap.Error(err))
			return err
		}
		if len(roles) != 0 {
			if err := grantUserRoles(tx, user.UID, roles); err != nil {
				zapx.WithContext(ctx).Error("grant user roles failed", zap.Error(err), zap.String("uid", user.UID), zap.Strings("roles", roles))
				return err
			}
		}
		return nil
	}); err != nil {
		return "", err
	}

	return user.UID, nil
}

func GetUserByUID(ctx context.Context, uid string) (*User, error) {
	user := &User{UID: uid}
	return getUser(ctx, user)
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{Email: email}
	return getUser(ctx, user)
}

func GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	user := &User{Phone: phone}
	return getUser(ctx, user)
}

func getUser(ctx context.Context, u *User) (*User, error) {
	db := core.DB.WithContext(ctx)
	if err := db.Where(u).First(u).Error; err != nil {
		zapx.WithContext(ctx).Error("get user from db failed", zap.Error(err), zap.Any("user", u))
		return nil, err
	}
	return u, nil
}

func UpdateUserInfo(ctx context.Context, u *User) (err error) {
	db := core.DB.WithContext(ctx)
	user := &User{}
	if err := db.Model(&User{UID: u.UID}).Where("uid = ?", u.UID).First(user).Error; err != nil {
		zapx.WithContext(ctx).Error("find user info failed", zap.Error(err), zap.String("UID", u.UID))
		return err
	}

	if u.Name != "" {
		user.Name = u.Name
	}

	if u.Email != "" {
		user.Email = u.Email
	}

	if u.Password != "" {
		u.Password, err = utils.EncryptPassword(u.Password)
		if err != nil {
			zapx.WithContext(ctx).Error("encrypt password failed", zap.Error(err))
			return err
		}
		user.Password = u.Password
	}

	if u.Gender != 0 {
		user.Gender = u.Gender
	}

	if u.AvatarURL != "" {
		user.AvatarURL = u.AvatarURL
	}

	if err := db.Model(&User{UID: u.UID}).Updates(user).Error; err != nil {
		zapx.WithContext(ctx).Error("update user info failed", zap.Error(err), zap.String("UID", u.UID))
		return err
	}
	return nil
}

func ResetPassword(ctx context.Context, phone, password string) error {
	zlog := zapx.WithContext(ctx)
	db := core.DB.WithContext(ctx)

	encryptedPassword, err := utils.EncryptPassword(password)
	if err != nil {
		zlog.With(zap.Error(err)).Error("encrypt password failed")
		return err
	}

	err = db.Model(&User{}).Where("phone = ?", phone).Update("password", encryptedPassword).Error
	if err != nil {
		zlog.With(zap.Error(err)).Error("update password error")
		return err
	}

	return nil
}

func GetUsersByUids(ctx context.Context, uids []string) ([]User, error) {
	var users []User
	db := core.DB.WithContext(ctx)
	if err := db.Where("uid in ?", uids).Find(&users).Error; err != nil {
		zapx.WithContext(ctx).Error("get users from db failed", zap.Error(err), zap.Any("user", uids))
		return nil, err
	}
	return users, nil
}

func GetGroupsDetail(ctx context.Context) (map[string]int, error) {
	db := core.DB.WithContext(ctx)
	var users []User
	if err := db.Select("groups").Where("groups is not NULL").Find(&users).Error; err != nil {
		zapx.WithContext(ctx).Error("get group details from db failed", zap.Error(err))
		return nil, err
	}

	groupsDetail := make(map[string]int)
	for _, user := range users {
		if len(user.Groups) != 0 {
			groupsDetail[user.Groups[0]] += 1
		}
	}
	return groupsDetail, nil
}
