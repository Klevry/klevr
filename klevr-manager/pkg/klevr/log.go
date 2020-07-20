package klevr

import (
	"strings"

	"github.com/NexClipper/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerEnv logger environment structure
type LoggerEnv struct {
	Level      string
	LogPath    string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// NewLoggerEnv constructor for default LoggerEnv
func NewLoggerEnv() *LoggerEnv {
	return &LoggerEnv{
		Level:      "debug",
		LogPath:    "./log/klevr.log",
		MaxSize:    20,
		MaxBackups: 20,
		MaxAge:     10,
		Compress:   false,
	}
}

// InitLogger init logger
func InitLogger(env *LoggerEnv) {
	setting := &lumberjack.Logger{
		Filename:   env.LogPath,
		MaxSize:    env.MaxSize,
		MaxBackups: env.MaxBackups,
		MaxAge:     5,
		Compress:   false,
	}

	logger.Init("StandardLogger", true, false, setting)

	var level logger.Level

	switch strings.ToLower(env.Level) {
	case "debug":
		level = 0
	case "info":
		level = 1
	case "warn", "warning":
		level = 2
	case "error":
		level = 3
	case "fatal":
		level = 4
	}

	logger.SetLevel(level)
	// logger 포맷 변경 (포맷 변경 시 파일:라인 정보 미표시로 주석처리)
	// logger.SetFlags(log.LstdFlags)
}
