package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	Dialect  string
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Charset  string
}

type DB struct {
	Db *sql.DB
}

func (DB *DB) Initialize(config DbConfig) (err error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Charset)
	if DB.Db, err = sql.Open(config.Dialect, dbURI); err != nil {
		return
	}
	if err = DB.Db.Ping(); err != nil {
		return
	}
	return
}
