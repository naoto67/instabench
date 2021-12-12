package lib

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type randomKeyRedisCtxKey struct{}

type RandomKeyRedisExecutor struct {
	client *redis.Client
	opt    RandomKeyRedisExecutorOption

	seed int64
}

type RandomKeyRedisExecutorOption struct {
	Host string
	Port string

	Clients int64

	KeyLen      int64
	DataSize    int64
	InitDataLen int64
}

func NewRandomKeyRedisExecutor(opt RandomKeyRedisExecutorOption) *RandomKeyRedisExecutor {
	return &RandomKeyRedisExecutor{
		seed: time.Now().Unix(),
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", opt.Host, opt.Port),
			PoolSize: int(opt.Clients),
		}),
		opt: opt,
	}
}

func (e *RandomKeyRedisExecutor) Setup(ctx context.Context) (context.Context, error) {
	d := generate(e.opt.DataSize)
	for i := 0; i < int(e.opt.InitDataLen); i++ {
		if err := e.client.Set(ctx, e.key(), d, 0).Err(); err != nil {
			return ctx, err
		}
	}
	time.Sleep(3 * time.Second)
	return ctx, nil
}

func (e *RandomKeyRedisExecutor) ExecContext(ctx context.Context) (context.Context, error) {
	v, ok := ctx.Value(&randomKeyRedisCtxKey{}).(string)
	if !ok {
		panic("not found context &randomKeyRedisCtxKey{}")
	}
	err := e.client.Get(ctx, v).Err()
	if errors.Is(err, redis.Nil) {
		return ctx, nil
	}
	return ctx, err
}

func (e *RandomKeyRedisExecutor) PrepareContext(ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, &randomKeyRedisCtxKey{}, e.key())
	return ctx, nil
}

func (e *RandomKeyRedisExecutor) key() string {
	return fmt.Sprintf("%0"+strconv.Itoa(int(e.opt.KeyLen))+"d", rand.Int63n(999999999))
}

func (e *RandomKeyRedisExecutor) ExtraConfig() map[string]interface{} {
	return map[string]interface{}{
		"host":          e.opt.Host,
		"port":          e.opt.Port,
		"clients":       e.opt.Clients,
		"init_data_len": e.opt.InitDataLen,
		"data_size":     e.opt.DataSize,
		"key_len":       e.opt.KeyLen,
	}
}
