package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexClipper/logger"

	"xorm.io/xorm/log"

	"xorm.io/xorm/names"

	_ "github.com/go-sql-driver/mysql" //justifying
	"xorm.io/xorm"
	// _ "github.com/jinzhu/gorm/dialects/sqlite" //justifying
	// // _ "github.com/jinzhu/gorm/dialects/mysql" //justfying
)

// DBInfo database connect info & connection info.
// URL: full URL for DB connect.
// MaxConnLifeTime: max connection life time(hour)
type DBInfo struct {
	Type            string
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifeTime int
	InitScriptPath  string
	ShowSql         bool
	LogLevel        string
	Cache           bool
}

type DB struct {
	*xorm.Engine
}

type Session struct {
	*xorm.Session
}

// Connect connect to Database and return DB
func (info *DBInfo) Connect() (*DB, error) {
	db, err := xorm.NewEngine(info.Type, info.URL)
	if err != nil {
		return &DB{db}, err
	}

	db.SetMaxOpenConns(info.MaxOpenConns)
	db.SetMaxIdleConns(info.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(info.MaxConnLifeTime) * time.Second)

	db.SetTableMapper(CustomTableNameMapper{})

	db.ShowSQL(info.ShowSql)

	switch strings.ToLower(info.LogLevel) {
	case "debug":
		db.SetLogLevel(log.LOG_DEBUG)
	case "info":
		db.SetLogLevel(log.LOG_INFO)
	case "warn":
		db.SetLogLevel(log.LOG_WARNING)
	case "error":
		db.SetLogLevel(log.LOG_ERR)
	default:
		db.SetLogLevel(log.LOG_OFF)
	}

	logger.Infof("DB log lever : [%s]", info.LogLevel)
	logger.Infof("DB show sql : [%v]", info.ShowSql)

	return &DB{db}, err
}

// NewSession New a session
func (db *DB) NewSession() *Session {
	return &Session{db.Engine.NewSession()}
}

func (tx *Session) Begin() error {
	return tx.Session.Begin()
}

type CustomTableNameMapper struct {
	base names.SnakeMapper
}

func (c CustomTableNameMapper) Obj2Table(s string) string {
	return strings.ToUpper(c.base.Obj2Table(s))
}

func (c CustomTableNameMapper) Table2Obj(s string) string {
	return c.base.Table2Obj(s)
}

// CheckGetQuery check query result
func CheckGetQuery(ok bool, err error) bool {
	if err != nil {
		panic(err)
	}

	return ok
}

// PanicForUpdate raise panic for update failure
func PanicForUpdate(t string, exeCnt int64, totCnt int64) {
	panic(fmt.Sprintf("All rows was not %s - (%d/%d)", t, exeCnt, totCnt))
}
