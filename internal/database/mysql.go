package database

import (
	"database/sql"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 连接池配置
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db

	slog.Info("database connected")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
