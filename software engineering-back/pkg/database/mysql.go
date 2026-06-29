package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	gsql "github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 通过 Connector 级别设置 Collation，确保每个新建的连接都使用 utf8mb4
	// 这比 SET NAMES 更可靠，因为 SET NAMES 只影响单个连接
	loc, _ := time.LoadLocation("Local")
	cfg := gsql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               dbName,
		ParseTime:            true,
		Loc:                  loc,
		Collation:            "utf8mb4_unicode_ci",
		AllowNativePasswords: true,
	}

	connector, err := gsql.NewConnector(&cfg)
	if err != nil {
		log.Fatalf("failed to create mysql connector: %v", err)
	}

	sqlDB := sql.OpenDB(connector)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB, err = gorm.Open(gormmysql.New(gormmysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                    true,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	log.Println("database connected successfully")
}
