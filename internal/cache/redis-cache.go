package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/naumovrus/finance-transaction-api/internal/entity"
	"github.com/naumovrus/finance-transaction-api/internal/service"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
	service *service.Service
}

func NewRedisCache(host string, db int, exp time.Duration, service1 *service.Service) TransactionCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: exp,
		service: service1,
	}
}

func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) SetTS(key string, post entity.TransactionSend) {
	client := cache.getClient()
	ctx := context.Background()
	// serialize Post object to JSON
	json, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	client.Set(ctx, key, json, cache.expires*time.Second)
}

func (cache *RedisCache) SetTB(key string, post *entity.TransactionTUTO) {
	client := cache.getClient()
	ctx := context.Background()
	// serialize Post object to JSON
	json, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	client.Set(ctx, key, json, cache.expires*time.Second)
}

func (cache *RedisCache) GetTS(key string) *entity.TransactionSend {
	client := cache.getClient()
	ctx := context.Background()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	transaction := entity.TransactionSend{}
	err = json.Unmarshal([]byte(val), &transaction)
	if err != nil {
		panic(err)
	}

	return &transaction
}

func (cache *RedisCache) GetTB(key string) *entity.TransactionTUTO {
	client := cache.getClient()
	ctx := context.Background()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	transaction := entity.TransactionTUTO{}
	err = json.Unmarshal([]byte(val), &transaction)
	if err != nil {
		panic(err)
	}

	return &transaction
}

func (cache *RedisCache) SetCachedData() error {
	client := cache.getClient()
	ctx := context.Background()
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		// * Scan all KEY, 20 each time
		keys, cursor, err = client.Scan(ctx, cursor, "*", 20).Result()
		if err != nil {
			log.Fatalf("%s ", err)
		}
		n += len(keys)

		log.Printf("\nfound %d keys\n", n)
		var value *entity.TransactionSend

		for _, key := range keys {
			value = cache.GetTS(key)
			log.Printf("%v %v %v %v\n", key, value.Id, value.UserIdFrom, value.Time)
			cache.service.Money.SetCachedDataSendPostgres(value.UserIdFrom, value.UserIdTo, value.Time)

		}

		pipe := client.Pipeline()
		for _, key := range keys {
			pipe.Del(ctx, key)
		}
		pipe.Exec(ctx)
		if cursor == 0 {
			break
		}
	}

	return nil
}
