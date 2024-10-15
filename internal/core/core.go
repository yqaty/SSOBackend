package core

import (
	"time"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/pkg/open"

	"github.com/UniqueStudio/UniqueSSOBackend/config"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/pkg/db"
	xredis "github.com/UniqueStudio/UniqueSSOBackend/internal/pkg/redis"
	"github.com/UniqueStudio/UniqueSSOBackend/internal/pkg/session"
	"github.com/gin-contrib/sessions"
	"github.com/go-redis/redis/v8"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client

	SessStore sessions.Store

	OpenClient *open.OpenClient
)

func Setup() (err error) {
	DB, err = db.NewPostgreConn(&db.Option{
		DSN:         config.Config.Database.Postgres.Dsn,
		MaxOpenConn: int(config.Config.Database.Postgres.MaxOpenConn),
		MaxIdleConn: int(config.Config.Database.Postgres.MaxIdleConn),
		MaxLifetime: time.Duration(config.Config.Database.Postgres.MaxLifeSeconds) * time.Second,
	})
	if err != nil {
		zapx.Error("new postgres connection failed", zap.Error(err))
		return err
	}

	RedisClient, err = xredis.NewRedisClient(&xredis.Option{DSN: config.Config.Database.Redis.Dsn})
	if err != nil {
		zapx.Error("new redis connection failed", zap.Error(err))
		return err
	}

	SessStore, err = session.NewSessionStore(&session.Option{
		RedisDSN:      config.Config.Database.Redis.Dsn,
		SessionSecret: config.Config.Application.SessionSecret,
		SessionDomain: config.Config.Application.SessionDomain,
	})
	if err != nil {
		zapx.Error("new session store failed", zap.Error(err))
		return err
	}

	//var conn *grpc.ClientConn
	//conn, err = grpc.Dial(
	//	config.Config.OpenPlatform.Addr,
	//	grpc.WithInsecure(),
	//	grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	//	grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	//)
	//if err != nil {
	//	zapx.Error("dial to open platform failed", zap.Error(err), zap.String("openPlatformAddr", config.Config.OpenPlatform.Addr))
	//	return err
	//}
	//SMSClient = open.NewSMSServiceClient(conn)
	//EmailClient = open.NewEmailServiceClient(conn)
	OpenClient = open.NewOpenClient(config.Config.OpenPlatform.Addr, config.Config.OpenPlatform.Token)

	if err != nil {
		zapx.Error("new lark client failed", zap.Error(err))
		return err
	}

	return nil
}
