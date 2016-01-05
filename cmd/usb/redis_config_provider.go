package main

import (
	"github.com/codegangsta/cli"
	"github.com/hpcloud/cf-usb/lib/config"
	"os"
	"strconv"
	"strings"

	"github.com/hpcloud/cf-usb/lib/config/redis"
)

type RedisConfigProvider struct {
}

func NewRedisConfigProvider() (*RedisConfigProvider, error) {
	return nil, nil
}

func (k *RedisConfigProvider) GetCLICommands(app Usb) []cli.Command {
	return []cli.Command{
		{
			Name: "redisConfigProvider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address, a",
					Usage: "Redis address and port",
				},
				cli.StringFlag{
					Name:  "database, d",
					Usage: "Redis database (integer)",
				},
				cli.StringFlag{
					Name:  "password,p",
					Usage: "Redis password",
				},
			},
			Action: redisConfigProviderCommand(app),
			Usage:  `Set redis configuration address`,
		},
	}
}

func redisConfigProviderCommand(app Usb) func(c *cli.Context) {
	return func(c *cli.Context) {
		logger := NewLogger(strings.ToLower(c.GlobalString("loglevel")))

		redisAddress := c.String("address")

		if redisAddress == "" {
			cli.ShowCommandHelp(c, "redisConfigProvider")
			os.Exit(0)
		}

		redisDatabase := c.String("database")
		redisPass := c.String("password")

		if redisDatabase != "" {
			db, err := strconv.ParseInt(redisDatabase, 10, 64)
			if err != nil {
				logger.Fatal("database must be a 64bit integer", err)
			}
			provisioner, err := redis.New(redisAddress, redisPass, db)
			if err != nil {
				logger.Fatal("redis config provider", err)
			}
			configuraiton := config.NewRedisConfig(provisioner)

			app.Run(configuraiton, logger)

		} else {
			provisioner, err := redis.New(redisAddress, redisPass, 0)
			if err != nil {
				logger.Fatal("redis config provider", err)
			}
			configuraiton := config.NewRedisConfig(provisioner)

			app.Run(configuraiton, logger)
		}

	}

}
