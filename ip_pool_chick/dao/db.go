package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"ip_pool_chick/config"
)

var DB *sql.DB

func InitDB(cfg config.DBConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	DB = db
	return nil
}
