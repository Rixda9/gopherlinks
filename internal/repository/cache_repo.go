package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository interface {
	Get(slug string) (string, error)
	Set(slug string, url string, expiry time.Duration) error
}

type RedisRepo struct {
	Client *redis.Client
	Context context.Context 
}

func NewRedisRepo(addr string) *RedisRepo {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: addr, 
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis.")

	return &RedisRepo{
		Client: rdb,
		Context: ctx,
	}
}

func (r *RedisRepo) Get(slug string) (string, error) {
	val, err := r.Client.Get(r.Context, slug).Result() 

	if err == redis.Nil {
		return "", nil 
	} else if err != nil {
		log.Printf("Redis GET error for slug %s: %v", slug, err)
		return "", fmt.Errorf("redis read error: %w", err)
	}
	
	return val, nil
}

func (r *RedisRepo) Set(slug string, url string, expiry time.Duration) error {
	err := r.Client.Set(r.Context, slug, url, expiry).Err()
	
	if err != nil {
		log.Printf("Redis SET error for slug %s: %v", slug, err)
		return fmt.Errorf("redis write error: %w", err)
	}
	
	return nil
}
