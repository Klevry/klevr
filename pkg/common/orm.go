package common

import (
	"time"

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

	return db, err
}

func init() {
	// // 기본 테이블명 설정 변경
	// // https://gorm.io/docs/conventions.html#Change-default-tablenames
	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return strings.ToUpper(defaultTableName)
	// }
}
