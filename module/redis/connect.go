package redis

import (
	"context"
	"errors"
	"github.com/bpcoder16/Water/logit"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func init() {
	loadRedisConfig()
}

func InitRedis() {
	connectDefault()
}

func connectDefault() {
	defaultRedis = redis.NewClient(&redis.Options{
		Addr:         config.Default.Host + ":" + strconv.Itoa(config.Default.Port),
		Username:     config.Default.Username,
		Password:     config.Default.Password,
		DB:           config.Default.DB,
		MaxRetries:   config.Default.MaxRetries,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
		PoolFIFO:     true,
		PoolSize:     100,
		//PoolTimeout:  200 * time.Millisecond,
		MinIdleConns:    20,
		MaxIdleConns:    50,
		ConnMaxIdleTime: 10 * time.Minute,
		//ConnMaxLifetime: 2 * time.Hour,
	})
	defaultRedis.AddHook(NewLoggerHook(logit.GetGlobalHelper()))
	err := defaultRedis.Get(context.Background(), "testConnect").Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		panic(config.Default.Host + ":" + strconv.Itoa(config.Default.Port) + ", failed to connect redis: " + err.Error())
	}
}
