package redis_service

import (
    "context"
	"time"
	"log"

    "github.com/redis/go-redis/v9"
    "github.com/bsm/redislock"
)

type RedisService struct {
    client *redis.Client
    locker *redislock.Client
}

func NewRedisService(host, password string, db int) *RedisService {
    // Create a Redis client instance
    client := redis.NewClient(&redis.Options{
        Addr:     host,
        Password: password,
        DB:       db,
    })

    // Ping the Redis server to test the connection
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        log.Fatalf("Failed to ping Redis server: %v", err)
    }

    log.Printf("Connected to Redis server")

    // Create a Redis lock client instance
    locker := redislock.New(client)

    // Return a new RedisService instance with the Redis client and lock client
    return &RedisService{client: client, locker: locker}
}

func (rs *RedisService) Close() error {
    // Close the Redis client connection
    return rs.client.Close()
}

func (rs *RedisService) AcquireLock(lockKey string, timeout int) (*redislock.Lock, error) {
    // Check out a Redis lock with the specified key and timeout
    lock, err := rs.locker.Obtain(context.Background(), lockKey, time.Duration(timeout)*time.Second, nil)
    if err != nil {
        log.Printf("[warning] failed acquiring Redis lock: %s", err)
        return nil, err
    }

    // Lock acquired successfully, return the lock
    return lock, nil
}