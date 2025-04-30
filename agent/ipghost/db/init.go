package db

import (
	"database/sql"
	"fmt"
	"ipghost/config"
	"time"
)

// InitDB 初始化数据库
func InitDB(conf config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		conf.Database.User,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 3)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	return db, nil
}
