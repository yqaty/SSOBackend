package session

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	goredis "github.com/go-redis/redis/v8"
)

type Option struct {
	RedisDSN      string
	SessionSecret    string
	SessionDomain string
}

func NewSessionStore(opt *Option) (sessions.Store, error) {
	rdsOpt, err := goredis.ParseURL(opt.RedisDSN)
	if err != nil {
		return nil, err
	}
	sessStore, err := redis.NewStoreWithDB(
		10, rdsOpt.Network, rdsOpt.Addr, rdsOpt.Password,
		strconv.FormatInt(int64(rdsOpt.DB), 10),
		[]byte(opt.SessionSecret),
	)
	if err != nil {
		return nil, err
	}
	sessStore.Options(sessions.Options{Path: "/", Domain: opt.SessionDomain, HttpOnly: true})
	return sessStore, nil
}
