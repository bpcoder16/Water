package nonblock

import (
	"context"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/module/redis"
	"strconv"
	"time"
)

func RedisLock(ctx context.Context, lockName string, deadLockExpireSecond time.Duration) bool {
	cacheValue := strconv.Itoa(int(time.Now().Add(deadLockExpireSecond).Unix()))
	success, err := redis.GetDefaultRedis().SetNX(ctx, lockName, cacheValue, deadLockExpireSecond).Result()

	if err != nil {
		logit.Context(ctx).WarnW("RedisLockErr", err.Error())
		return false
	}

	// 防止死锁
	if !success {
		if expireTimeStr, errRedis := redis.GetDefaultRedis().Get(ctx, lockName).Result(); errRedis == nil {
			if expireTimeRedis, errStr := strconv.Atoi(expireTimeStr); errStr == nil {
				if time.Now().Unix() > int64(expireTimeRedis) {
					redis.GetDefaultRedis().Del(ctx, lockName)
				}
			} else {
				redis.GetDefaultRedis().Del(ctx, lockName)
			}
		} else {
			redis.GetDefaultRedis().Del(ctx, lockName)
		}
	}

	return success
}

func RedisUnlock(ctx context.Context, lockName string) {
	redis.GetDefaultRedis().Del(ctx, lockName)
}
