package common

import (
	"net/http"
	"strings"
	"time"

	"github.com/NexClipper/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

// LoggerEnv logger environment structure
type LoggerEnv struct {
	Level      string `envconfig:"LOG_LEVEL"`
	LogPath    string `envconfig:"LOG_PATH"`
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

type LogFormatter func(params LogFormatterParams) string

type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return true
}

func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}
