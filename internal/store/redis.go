package store

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type RedisStore struct {
    Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
    return &RedisStore{
        Client: redis.NewClient(&redis.Options{
            Addr: addr,
        }),
    }
}

func (rs *RedisStore) SetOTP(ctx context.Context, phone string, otp string, ttl time.Duration) error {
    return rs.Client.Set(ctx, "otp:"+phone, otp, ttl).Err()
}

func (rs *RedisStore) GetOTP(ctx context.Context, phone string) (string, error) {
    return rs.Client.Get(ctx, "otp:"+phone).Result()
}

func (rs *RedisStore) DeleteOTP(ctx context.Context, phone string) error {
    return rs.Client.Del(ctx, "otp:"+phone).Err()
}

// RateLimitSend increments and returns current count. It sets TTL only when key is new (i.e., value==1).
func (rs *RedisStore) IncSendCount(ctx context.Context, phone string, window time.Duration) (int64, error) {
    key := "rate:send:" + phone
    n, err := rs.Client.Incr(ctx, key).Result()
    if err != nil {
        return 0, err
    }
    if n == 1 {
        // set expiry on first increment
        rs.Client.Expire(ctx, key, window)
    }
    return n, nil
}

func (rs *RedisStore) GetSendCount(ctx context.Context, phone string) (int64, error) {
    key := "rate:send:" + phone
    return rs.Client.Get(ctx, key).Int64()
}

// Attempts
func (rs *RedisStore) IncAttempt(ctx context.Context, phone string, ttl time.Duration) (int64, error) {
    key := "attempts:" + phone
    n, err := rs.Client.Incr(ctx, key).Result()
    if err != nil {
        return 0, err
    }
    if n == 1 {
        rs.Client.Expire(ctx, key, ttl)
    }
    return n, nil
}

func (rs *RedisStore) ResetAttempts(ctx context.Context, phone string) error {
    return rs.Client.Del(ctx, "attempts:"+phone).Err()
}

// Block
func (rs *RedisStore) SetBlocked(ctx context.Context, phone string, ttl time.Duration) error {
    return rs.Client.Set(ctx, "blocked:"+phone, "1", ttl).Err()
}
func (rs *RedisStore) IsBlocked(ctx context.Context, phone string) (bool, error) {
    ok, err := rs.Client.Exists(ctx, "blocked:"+phone).Result()
    if err != nil { return false, err }
    return ok == 1, nil
}
