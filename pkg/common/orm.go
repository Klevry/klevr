package common

import (
	"fmt"
	"strings"
	"time"

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
}

// Connect connect to Database and return DB
func (info *DBInfo) Connect() (*xorm.Engine, error) {
	db, err := xorm.NewEngine(info.Type, info.URL)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(info.MaxOpenConns)
	db.SetMaxIdleConns(info.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(info.MaxConnLifeTime) * time.Second)

	db.SetTableMapper(CustomTableNameMapper{})

	return db, err
}

func init() {
	// // 기본 테이블명 설정 변경
	// // https://gorm.io/docs/conventions.html#Change-default-tablenames
	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return strings.ToUpper(defaultTableName)
	// }
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
