package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Option struct {
	DSN string
}

func NewRedisClient(opt *Option) (*redis.Client, error) {
	rdsOpt, err := redis.ParseURL(opt.DSN)
	if err != nil {
		return nil, err
	}

	rclient := redis.NewClient(rdsOpt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = rclient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rclient, nil
}
