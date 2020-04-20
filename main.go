package main

import (
	"os"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/urfave/cli"

	dconfig "github.com/canyanio/rating-agent-janus/config"
	"github.com/canyanio/rating-agent-janus/server"
)

func main() {
	doMain(os.Args)
}

func doMain(args []string) {
	var configPath string
	var configDebug bool

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Configuration `FILE`. Supports JSON, TOML, YAML and HCL formatted configs.",
				Value:       "config.yaml",
				Destination: &configPath,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug mode and verbose logging",
				Destination: &configDebug,
			},
		},
		Commands: []cli.Command{
			{
				Name:   "agent",
				Usage:  "Run the Janus agent",
				Action: cmdAgent,
				Flags:  []cli.Flag{},
			},
		},
	}
	app.Usage = "rating-agent-janus"
	app.Version = "1.0.0"
	app.Action = cmdAgent

	app.Before = func(args *cli.Context) error {
		err := dconfig.Init(configPath)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		log.Setup(configDebug)

		return nil
	}

	err := app.Run(args)
	if err != nil {
		cli.NewExitError(err.Error(), 1)
	}
}

func cmdAgent(args *cli.Context) error {
	return server.InitAndRun(config.Config)
}
