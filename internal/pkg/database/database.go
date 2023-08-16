package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
)

var (
	db     *sqlx.DB
	config DBConfig   // 数据库链接配置
	once   sync.Mutex // 数据库链接初始化锁
)

// DBConfig 数据库链接配置
type DBConfig struct {
	dsn          string
	driver       string
	maxOpenConns int
	maxIdleConns int
}

func NewDBConfig(dsn, driver string, maxOpenConns, maxIdleConns int) {
	config = DBConfig{
		dsn:          dsn,
		driver:       driver,
		maxOpenConns: maxOpenConns,
		maxIdleConns: maxIdleConns,
	}
}

// 连接数据库
func connect() {
	once.Lock()
	defer once.Unlock()
	// 避免排队线程重复初始化
	if db == nil {
		var err error
		db, err = sqlx.Open(config.driver, config.dsn)
		if err != nil {
			panic(fmt.Sprintf("数据库连接失败。原因：%s\n", err.Error()))
		}
		db.SetMaxOpenConns(config.maxOpenConns)
		db.SetMaxIdleConns(config.maxIdleConns)
	}
}

// GetDBInstance 获取数据库链接实例
func GetDBInstance() *sqlx.DB {
	if db == nil {
		connect()
	}
	return db
}

func Run(txFunc func(db *sqlx.Tx) error) error {
	db := GetDBInstance()
	tx, err := db.BeginTxx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("开启数据库事务失败。%s", err.Error())
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

type DB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)

	Exec(query string, args ...any) (sql.Result, error)

	Get(dest interface{}, query string, args ...interface{}) error

	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}
