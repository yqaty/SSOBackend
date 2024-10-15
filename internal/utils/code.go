package utils

import (
	"context"
	"time"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/core"
	"github.com/labstack/gommon/random"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

func newTmpCode() string {
	return random.String(6, random.Numeric)
}

func GenerateTmpCode(ctx context.Context, id string, expire time.Duration) (string, error) {
	code := newTmpCode()
	err := core.RedisClient.Set(ctx, id, code, expire).Err()
	if err != nil {
		zapx.WithContext(ctx).Error("generate tmp code failed", zap.Error(err), zap.String("id", id))
		return "", err
	}
	return code, nil
}

// GetTmpCodeByID - it will delete the validation code whether it validate
func GetTmpCodeByID(ctx context.Context, id string) (code string, err error) {
	value := core.RedisClient.GetDel(ctx, id)
	if err = value.Err(); err != nil {
		zapx.WithContext(ctx).Error("getdel by id failed", zap.Error(err), zap.String("id", id))
		return "", err
	}
	if err = value.Scan(&code); err != nil {
		return
	}
	return
}
