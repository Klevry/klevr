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

	return config, nil
}

const (
	c, ll, lp, p, hook string = "config", "log.level", "log.path", "port", "webhook.url"
)

func main() {
	// TimeZone UTC로 설정
	os.Setenv("TZ", "")

	var exit int = 0

	app := &cli.App{
		Name:      "Klevr-Manager",
		Version:   "v1.0.0",
		Copyright: "(c) 2020 NexCloud",
		Usage:     "main [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     c,
				Aliases:  []string{"c"},
				Value:    "./conf/klevr-manager-local.yml",
				Usage:    "Config file path",
				Required: true,
				EnvVars:  []string{"KLEVR_CONFIG_PATH"},
			},
			&cli.StringFlag{
				Name:     ll,
				Aliases:  []string{"ll"},
				Value:    "debug",
				Usage:    "Logging level(default:debug, info, warn, error, fatal)",
				Required: false,
				EnvVars:  []string{"LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:     lp,
				Aliases:  []string{"lp"},
				Value:    "./log/klevr-manager.log",
				Usage:    "log full path(include file name)",
				Required: false,
				EnvVars:  []string{"LOG_PATH"},
			},
			&cli.StringFlag{
				Name:     p,
				Aliases:  []string{"p"},
				Value:    "8090",
				Usage:    "default port used by the klevr-manager(default:8090)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     hook,
				Aliases:  []string{"hook"},
				Usage:    "WebHook URL",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			// 설정파일 반영
			config, err := loadConfig(ctx.String(c))
			if err != nil {
				logger.Fatal(err)
				exit = 1

				panic("Can not start klevr-manager")
			}

			// 환경변수 반영
			envconfig.Process("", config)
			envAssembledConfig := *config

			// 실행 파라미터 반영 (실행 파라미터>환경변수>설정파일)
			if ctx.String(ll) != "" {
				config.Log.Level = ctx.String(ll)
			}
			if ctx.String(lp) != "" {
				config.Log.LogPath = ctx.String(lp)
			}
			if ctx.String(p) != "" {
				config.Klevr.Server.Port = ctx.Int(p)
			}
			if ctx.String(hook) != "" {
				config.Klevr.Server.Webhook.Url = ctx.String(hook)
			}

			loggerEnv := &common.LoggerEnv{
				Level:      config.Log.Level,
				LogPath:    config.Log.LogPath,
				MaxSize:    config.Log.MaxSize,
				MaxBackups: config.Log.MaxBackups,
				MaxAge:     config.Log.MaxAge,
				Compress:   config.Log.Compress,
			}
			common.InitLogger(loggerEnv)

			logger.Info("Start Klevr-manager")
			logger.Debug("ENV assembled config : ", &envAssembledConfig)

			// Actual instance running point
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
}
