package pubsub

import (
	"context"
	wRedis "github.com/bpcoder16/Water/module/redis"
	"github.com/redis/go-redis/v9"
)

type RedisPubSub struct {
	channels []string
}

func NewRedisPubSub(channels ...string) *RedisPubSub {
	return &RedisPubSub{
		channels: channels,
	}
}

func (r *RedisPubSub) Subscribe(f func(*redis.Message)) error {
	ctx := context.Background()
	pubSub := wRedis.GetDefaultRedis().Subscribe(ctx, r.channels...)
	defer func() {
		_ = pubSub.Close()
	}()

	for {
		msg, errR := pubSub.ReceiveMessage(ctx)
		if errR != nil {
			return errR
		}
		f(msg)
	}
}
