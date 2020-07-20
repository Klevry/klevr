package klevr

import (
	"strings"

	_ "github.com/go-sql-driver/mysql" //justfying
	"github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/mysql" //justfying
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
func (info *DBInfo) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(info.Type, info.URL)
	if err != nil {
		return db, err
	}

	return db, err
}

func init() {
	// 기본 테이블명 설정 변경
	// https://gorm.io/docs/conventions.html#Change-default-tablenames
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return strings.ToUpper(defaultTableName)
	}
}
