package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"ip_mole/config"
)

func InitDB(conf config.Config) (*sql.DB, error) {
	var cfg = conf.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
