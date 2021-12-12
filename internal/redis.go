package internal

import (
	"github.com/naoto67/instabench"
	"github.com/naoto67/instabench/lib"
	"github.com/urfave/cli/v2"
)

func redisAction(cliContext *cli.Context) error {
	var executor instabench.Executor
	executor = lib.NewRandomKeyRedisExecutor(lib.RandomKeyRedisExecutorOption{
		Host:        redisHost,
		Port:        redisPort,
		Clients:     redisClients,
		KeyLen:      6,
		DataSize:    100,
		InitDataLen: 100000,
	})
	return exec(cliContext, executor)
}
