package main

import (
	"io/ioutil"
	"os"
	"sort"

	"github.com/NexClipper/logger"
	"github.com/kelseyhightower/envconfig"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/manager"
	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"
)

type config struct {
	Log   common.LoggerEnv
	Klevr manager.Config
}

func loadConfig(configPath string) (*config, error) {
	logger.Debug("configPath : ", configPath)
	var err error

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, common.NewStandardErrorWrap("configuration loading failed", err)
	}

	config := &config{}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, common.NewStandardErrorWrap("configuration loading failed", err)
	}

	logger.Debug("loaded config : ", *config)

	return config, nil
}

func main() {
	common.InitLogger(common.NewLoggerEnv())

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
				Required: true,
				EnvVars:  []string{"KLEVR_CONFIG_PATH"},
			},
			&cli.StringFlag{
				Name:     "log.level",
				Aliases:  []string{"ll"},
				Value:    "debug",
				Usage:    "Logging level(default:debug, info, warn, error, fatal)",
				Required: false,
				EnvVars:  []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:     "log.path",
				Aliases:  []string{"lp"},
				Value:    "./log/klevr-manager.log",
				Usage:    "log full path(include file name)",
				Required: false,
				EnvVars:  []string{"LOG_PATH"},
			},
			&cli.StringFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Value:    "debug",
				Usage:    "Logging level(default:debug, info, warn, error, fatal)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "webhook.url",
				Aliases:  []string{"hook"},
				Usage:    "WebHook URL",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			// 설정파일 반영
			config, err := loadConfig(c.String("config"))
			if err != nil {
				logger.Fatal(err)
				exit = 1

				panic("Can not start klevr-manager")
			}

			// 환경변수 반영
			envconfig.Process("", config)

			logger.Debug("ENV assembled config : ", *config)

			// 실행 파라미터 반영 (실행 파라미터>환경변수>설정파일)
			if c.String("log.level") != "" {
				config.Log.Level = c.String("log.level")
			}
			if c.String("log.path") != "" {
				config.Log.LogPath = c.String("log.path")
			}
			if c.String("port") != "" {
				config.Klevr.Server.Port = c.Int("port")
			}
			if c.String("webhook.url") != "" {
				config.Klevr.Server.Webhook.Url = c.String("webhook.url")
			}

			// common.ContextPut("appConfig", config)
			// common.ContextPut("cliContext", c)

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
