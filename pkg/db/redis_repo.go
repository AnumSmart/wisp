package db

import (
	"context"
	"log"
	"simple_gin_server/configs"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReddisRepoInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Exists(ctx context.Context, redisKey string) (bool, error)
}

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(ctx context.Context, conf *configs.Config) *RedisRepo {
	reddis_db_n, err := strconv.Atoi(conf.Redis.NDB)
	if err != nil {
		log.Println("Could not read reddis DB number from config")
	}
	return &RedisRepo{
		redis.NewClient(
			&redis.Options{
				Addr:     conf.Redis.Addr,
				Password: conf.Redis.Pass,
				DB:       reddis_db_n,
			},
		),
	}
}

func (r *RedisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisRepo) Exists(ctx context.Context, redisKey string) (bool, error) {
	result, err := r.client.Exists(ctx, redisKey).Result()
	if err != nil {
		log.Printf("[redis_repo.go]---[Exists()], Err:%v", err)
		return false, err
	}

	// Возвращает true, если ключ существует (1), false если нет (0)
	return result == 1, nil
}
