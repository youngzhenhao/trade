package services

import (
	"context"
	"errors"
	"math/rand"
	"time"
	"trade/middleware"
	"trade/utils"
)

func RedisSetRand(name string, duration time.Duration) (string, error) {
	redis := middleware.Client
	ctx := context.Background()
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	captcha := utils.RandString(seed, 64)
	err := redis.Set(ctx, name+"_captcha", captcha, duration).Err()
	return captcha, err
}

func RedisGetVerify(name string, captcha string) error {
	if captcha == "" {
		return errors.New("captcha is empty")
	}
	redis := middleware.Client
	ctx := context.Background()
	result, err := redis.Get(ctx, name+"_captcha").Result()
	if err != nil {
		return err
	}
	if result != captcha {
		return errors.New("captcha is wrong")
	}
	return nil
}
