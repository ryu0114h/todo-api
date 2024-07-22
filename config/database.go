package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB は、データベース接続を行う。
// 最大10回の試行を行い、2秒間隔で再試行する。
// 接続が成功した場合、データベースのロガーを Info レベルに設定し、データベースインスタンスを返す。
// 10回の試行後に接続できない場合、Exitする。
func InitDB() *gorm.DB {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		slog.Info(fmt.Sprintf("failed to connect database, retrying in 2 seconds: %v", err))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil
	}
	db.Logger = db.Logger.LogMode(logger.Info)
	return db
}
