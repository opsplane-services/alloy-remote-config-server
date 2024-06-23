package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type Storage struct {
	UseRedis     bool
	MemoryStore  map[string]string
	RedisStore   *redis.Client
	Organization string
	TTL          time.Duration
}

func InitStorage() (*Storage, error) {
	orgName := os.Getenv("ORG_NAME")
	useRedis := os.Getenv("USE_REDIS") == "true"
	if useRedis {
		redisOpt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
		if err != nil {
			return nil, err
		}
		redisOpt.DB = 0
		redisOpt.IdleTimeout = time.Second * 60
		redisOpt.IdleCheckFrequency = time.Second * 5
		redisClient := redis.NewClient(redisOpt)
		ttlNum := 259200
		ttlEnv := os.Getenv("REDIS_TTL")
		if len(ttlEnv) > 0 {
			ttlNum, _ = strconv.Atoi(ttlEnv)
		}
		ttlDuration := time.Duration(ttlNum) * time.Second
		return &Storage{RedisStore: redisClient, UseRedis: true, Organization: orgName, TTL: ttlDuration}, nil
	} else {
		return &Storage{MemoryStore: make(map[string]string), UseRedis: false, Organization: orgName}, nil
	}
}

func (s *Storage) Set(id string, content string) error {
	key := fmt.Sprintf("{%s}:%s", s.Organization, id)
	if s.UseRedis {
		ctx := context.Background()
		err := s.RedisStore.Set(ctx, key, content, 0).Err()
		if err != nil {
			return err
		}
		_, err = s.RedisStore.Expire(ctx, key, s.TTL).Result()
		if err != nil {
			return err
		}
	} else {
		s.MemoryStore[key] = content
	}
	return nil
}

func (s *Storage) Get(id string) (string, error) {
	key := fmt.Sprintf("{%s}:%s", s.Organization, id)
	if s.UseRedis {
		ctx := context.Background()
		value, err := s.RedisStore.Get(ctx, key).Result()
		if err != nil {
			return "", fmt.Errorf("Key (id) does not exist: %s", id)
		}
		return value, nil
	} else {
		value, exists := s.MemoryStore[key]
		if !exists {
			return "", fmt.Errorf("Key (id) does not exist: %s", id)
		}
		return value, nil
	}
}

func (s *Storage) GetAll() ([]string, error) {
	result := make([]string, 0)
	orgPrefix := fmt.Sprintf("{%s}:", s.Organization)
	if s.UseRedis {
		ctx := context.Background()
		scanCount := 100
		match := fmt.Sprintf("%s*", orgPrefix)
		var cursor uint64
		for {
			var ks []string
			var err error
			ks, cursor, err = s.RedisStore.Scan(ctx, cursor, match, int64(scanCount)).Result()
			if err != nil {
				return nil, err
			}
			for _, k := range ks {
				result = append(result, strings.TrimPrefix(k, orgPrefix))
			}
			if cursor == 0 {
				break
			}
		}
	} else {
		for k, _ := range s.MemoryStore {
			result = append(result, strings.TrimPrefix(k, orgPrefix))
		}
	}
	return result, nil
}
