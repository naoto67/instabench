package main

import (
	"bufio"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/naoto67/instabench/internal"
)

type RedisCommand struct {
	client *redis.Client
}

func NewRedisCommand() *RedisCommand {
	return &RedisCommand{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			PoolSize: 2000,
		}),
	}
}

func (r *RedisCommand) ExecContext(ctx context.Context) (result *Result, err error) {
	began := time.Now()
	result = &Result{Timestamp: began}
	defer func() {
		result.Latency = time.Since(began)
	}()
	conn := internal.Get()
	defer internal.Put(conn)
	_, err = conn.Write([]byte("GET KEY\r\n"))
	if err != nil {
		return result, err
	}
	_, err = bufio.NewReader(conn).ReadString('\n')
	// err = r.client.Get(ctx, "key").Err()
	return result, err
}
