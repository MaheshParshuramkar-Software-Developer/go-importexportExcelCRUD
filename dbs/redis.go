package dbs

import (
	"context"
	"time"

	"github.com/go-importexportExcelCRUD/logger"
	"github.com/go-importexportExcelCRUD/models"
)

var RClient *RedisClient
var redisConf *models.RedisConfig

// RedisClient : redis client which stores client and flag whether the connection is active or not
type RedisClient struct {
	Client      *redis.Client
	IsConnected bool
}

// NewRedisClient : Creating new redis client and returning it
func NewRedisClient(config *models.RedisConfig) *RedisClient {
	redisCli := new(RedisClient)
	redisConf = config
	redisCli.connect()
	if !redisCli.IsConnected {
		logger.Log.Fatal("error while connecting to  redis")
	}
	RClient = redisCli
	return redisCli
}

func (redisCli *RedisClient) connect() {
	addr := redisConf.Host + ":" + redisConf.Port
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConf.Password, // provide password of redis server
		DB:       redisConf.RedisDB,  // use database 1
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		logger.Log.Debug("error while ping to redis ", err.Error())
		redisCli.IsConnected = false
		return
	}
	redisCli.IsConnected = true
	redisCli.Client = client
}

// Get the redis key
func (redisCli *RedisClient) Get(key string) (string, error) {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.Get(context.Background(), key).Result()
}

// Set the redis key with value
func (redisCli *RedisClient) Set(key string, value string, expiry time.Duration) error {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.Set(context.Background(), key, value, expiry).Err()
}

// Del the redis key
func (redisCli *RedisClient) Del(key string) error {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.Del(context.Background(), key).Err()
}

// LPush store the values at the front of queue associated with a specific key
func (redisCli *RedisClient) LPush(key string, values ...interface{}) error {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.LPush(context.Background(), key, values).Err()
}

// RPush store the values at the end of queue associated with a specific key
func (redisCli *RedisClient) RPush(key string, values ...interface{}) error {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.RPush(context.Background(), key, values).Err()
}

// LPop pop a element from the front of queue associated with a specific key
func (redisCli *RedisClient) LPop(key string) (string, error) {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.LPop(context.Background(), key).Result()
}

// LPopCount pop no of elements from the front of queue associated with a specific key
func (redisCli *RedisClient) LPopCount(key string, count int) ([]string, error) {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.LPopCount(context.Background(), key, count).Result()
}

// RPop pop a element from the end of queue associated with a specific key
func (redisCli *RedisClient) RPop(key string) (string, error) {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.RPop(context.Background(), key).Result()
}

// LRange get the values in specific range from queue associated with a specific key
func (redisCli *RedisClient) LRange(key string, start int64, end int64) ([]string, error) {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.LRange(context.Background(), key, start, end).Result()
}

// LLen to get the length of the queue associated with a specific key
func (redisCli *RedisClient) LLen(key string) (int64, error) {
	if !redisCli.IsConnected {
		redisCli.connect()
	}
	return redisCli.Client.LLen(context.Background(), key).Result()
}

// Close the redis connection whenever the job is done or shutdown is received
func (redisCli *RedisClient) Close() {
	redisCli.IsConnected = false
	if err := redisCli.Client.Close(); err != nil {
		logger.Log.Debug("error while closing redis client", err.Error())
	}
	logger.Log.Debug("redis client closed ...")
}
