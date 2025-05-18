package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/i474232898/chatserver/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pool    *gorm.DB
	once    sync.Once
	poolErr error
)

func Connect(cfg *configs.AppConfigs) error {
	once.Do(func() {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: "user=postgres host=localhost password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		}), &gorm.Config{})
		pool = db
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			poolErr = err
			return
		}
	})

	return poolErr
}

func GetPool() *gorm.DB {
	return pool
}
