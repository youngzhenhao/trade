package middleware

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"sync"
	"time"
	"trade/config"
)

var (
	ctx      = context.Background()
	Client   *redis.Client
	ridesOne sync.Once
)

func RedisConnect() error {
	var errs error
	once.Do(func() {
		loadConfig, err := config.LoadConfig("config.yaml")
		if err != nil {
			panic("failed to load config: " + err.Error())
		}
		redisAddr := fmt.Sprintf("%s:%s", loadConfig.Redis.Host, loadConfig.Redis.Port)
		Client = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Username: loadConfig.Redis.Username,
			Password: loadConfig.Redis.Password,
			DB:       loadConfig.Redis.DB,
			PoolSize: 2000,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// 检查连接是否成功
		if _, err := Client.Ping(ctx).Result(); err != nil {
			errs = fmt.Errorf("failed to connect to Redis: %w", err)
		}
	})
	return errs
}

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

func RedisGet(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

func RedisDel(key string) error {
	return Client.Del(ctx, key).Err()
}

func AcquireLock(key string, expiration time.Duration) (string, bool) {
	uniqueID := uuid.New().String() // Generate a unique identifier
	result, err := Client.SetNX(ctx, key, uniqueID, expiration).Result()
	if err != nil {
		log.Println("Error acquiring lock:", err)
		return "", false
	}
	if result {
		return uniqueID, true
	}
	return "", false
}

func ReleaseLock(key string, identifier string) {
	luaScript := `
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
    `
	result, err := Client.Eval(ctx, luaScript, []string{key}, identifier).Result()
	if err != nil {
		log.Println("Error releasing lock:", err)
	} else if result.(int64) != 1 {
		log.Println("Failed to release lock, not the owner")
	}
}
