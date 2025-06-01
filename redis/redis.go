package redis

import (
	"time"

	"github.com/danielpnjt/go-library/log"
	"github.com/go-redis/redis"
)

type RedisOop struct {
	redisClient *redis.Client
}

const (
	Nil = redis.Nil
)

func Init(url string, password string) (*RedisOop, error) {
	redClient := &RedisOop{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: password,
		}),
	}

	return redClient, nil
}

func (r *RedisOop) SetRedisString(key string, otp string, expiration time.Duration) error {
	err := r.redisClient.Set(key, otp, expiration).Err()
	if err != nil {
		log.LogDebug("Error SetRedisString: " + err.Error())
	}
	return err
}

func (r *RedisOop) Get(key string) (string, error) {
	attemptString, err := r.redisClient.Get(key).Result()
	if err != nil {
		log.LogDebug("Error GetRedis: " + err.Error())
	}
	return attemptString, err
}

func (r *RedisOop) SetRedisHash(key string, objectRedis map[string]interface{}, expiration time.Duration) error {
	err := r.redisClient.HMSet(key, objectRedis).Err()
	r.redisClient.Expire(key, expiration)
	if err != nil {
		log.LogDebug("Error SetRedisHash: " + err.Error())
	}
	return err
}

func (r *RedisOop) Increase(key string, field string) error {
	err := r.redisClient.HIncrBy(key, field, 1).Err()
	if err != nil {
		log.LogDebug("Error Increase: " + err.Error())
	}
	return err
}

func (r *RedisOop) Delete(key string) error {
	err := r.redisClient.Del(key).Err()
	if err != nil {
		log.LogDebug("Error Delete: " + err.Error())
	}
	return err
}

func (r *RedisOop) GetHash(key string) (map[string]string, error) {
	data, err := r.redisClient.HGetAll(key).Result()
	if err != nil {
		log.LogDebug("Error GetHash: " + err.Error())
	}

	return data, err
}

func (r *RedisOop) GetTTLInSecond(key string) (int, error) {
	cd, err := r.redisClient.TTL(key).Result()
	inSecond := int(cd.Seconds())

	if err != nil {
		log.LogDebug("Error GetTTLInSecond: " + err.Error())
	}
	return inSecond, err
}

func (r *RedisOop) IncreaseByKey(key string) (int64, error) {
	num, err := r.redisClient.Incr(key).Result()
	if err != nil {
		log.LogDebug("Error IncreaseByKey: " + err.Error())
	}

	return num, err
}
