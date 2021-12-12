package internal

import (
	"errors"

	"github.com/urfave/cli/v2"
)

func NewBenchCommand() *cli.Command {
	return &cli.Command{
		Name: "bench",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        serviceFlag,
				Destination: &service,
				Value:       "redis",
			},
			&cli.Int64Flag{
				Name:        threadFlag,
				Destination: &thread,
				Value:       1000,
			},
			&cli.Int64Flag{
				Name:        durationFlag,
				Usage:       "sec",
				Value:       2,
				Destination: &duration,
			},
			&cli.StringFlag{
				Name: outputJsonPath,
			},

			// redis option
			&cli.StringFlag{
				Name:        redisHostFlag,
				Destination: &redisHost,
				Value:       "localhost",
			},
			&cli.StringFlag{
				Name:        redisPortFlag,
				Destination: &redisPort,
				Value:       "6379",
			},
			&cli.Int64Flag{
				Name:        redisClientsFlag,
				Destination: &redisClients,
				Value:       1000,
			},
		},
		Action: func(cliContext *cli.Context) error {
			switch service {
			case "redis":
				redisAction(cliContext)
			default:
				return errors.New("not supported service")
			}
			return nil
		},
	}
}
