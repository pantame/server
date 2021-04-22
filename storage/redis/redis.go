package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/pantame/server/config"
	"time"
)

type Storage struct {
	db *redis.Client
}

func New(conf config.RedisConfig) *Storage {
	db := redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		Password: conf.Password,
		DB:       conf.DB,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}
}

func (s *Storage) Get(key string) (string, error) {
	if len(key) <= 0 {
		return "", errors.New("storage: invalid parameter")
	}

	val, err := s.db.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", errors.New("storage: not found")
	}
	return val, err
}

func (s *Storage) GetInt(key string) (int, error) {
	if len(key) <= 0 {
		return 0, errors.New("storage: invalid parameter")
	}

	val, err := s.db.Get(context.Background(), key).Int()
	if err == redis.Nil {
		return 0, errors.New("storage: not found")
	}
	return val, err
}

func (s *Storage) GetBytes(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, errors.New("storage: invalid parameter")
	}

	val, err := s.db.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, errors.New("storage: not found")
	}
	return val, err
}

func (s *Storage) Set(key string, value interface{}, exp time.Duration) error {
	if len(key) <= 0 {
		return errors.New("storage: invalid parameter")
	}
	return s.db.Set(context.Background(), key, value, exp).Err()
}

func (s *Storage) IncrBy(key string, val int64) error {
	if len(key) <= 0 {
		return errors.New("storage: invalid parameter")
	}
	return s.db.IncrBy(context.Background(), key, val).Err()
}

func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return errors.New("storage: invalid parameter")
	}
	return s.db.Del(context.Background(), key).Err()
}

func (s *Storage) Reset() error {
	return s.db.FlushDB(context.Background()).Err()
}

func (s *Storage) Close() error {
	return s.db.Close()
}
