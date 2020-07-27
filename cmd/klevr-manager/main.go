package main

import (
	"io/ioutil"
	"os"
	"sort"

	"github.com/NexClipper/logger"

	klevr "github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/manager"
	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"
)

type config struct {
	Log   klevr.LoggerEnv
	Klevr manager.Config
}

func loadConfig(configPath string) (*config, error) {
	logger.Debug("configPath : ", configPath)
	var err error

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, klevr.NewCheckedErrorWrap("configuration loading failed", &err)
	}

	config := &config{}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, klevr.NewCheckedErrorWrap("configuration loading failed", &err)
	}

	logger.Debug("loaded config : ", *config)

	return config, nil
}

func main() {
	klevr.InitLogger(klevr.NewLoggerEnv())

	logger.Info("Start Klevr-manager")

	var exit int = 0

	app := &cli.App{
		Name:      "Klevr-Manager",
		Version:   "v1.0.0",
		Copyright: "(c) 2020 NexCloud",
		Usage:     "main [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Value:    "./conf/klevr-manager-local.yml",
				Usage:    "Config file path",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "log.level",
				Aliases:  []string{"L"},
				Value:    "debug",
				Usage:    "Logging level(default:debug, info, warn, error, fatal)",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			config, err := loadConfig(c.String("config"))
			if err != nil {
				logger.Fatal(err)
				exit = 1

				panic("Can not start klevr-manager")
			}

			if c.String("log.level") != "" {
				config.Log.Level = c.String("log.level")
			}

			/// Actual instance running point
			instance, err := manager.NewKlevrManager()
			if err != nil {
				logger.Error(err)
			}

			instance.SetConfig(&config.Klevr)
			instance.Run()

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
	}

	defer logger.Info("Stopped Klevr-manager")
	defer logger.Close()
	defer os.Exit(exit)

	//os.Exit(run())
}
