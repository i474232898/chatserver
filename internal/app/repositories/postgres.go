package repositories

import (
	"fmt"
	"log/slog"
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
			DSN: fmt.Sprintf(
				"user=%s host=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				cfg.DB.User,
				cfg.DB.Host,
				cfg.DB.Password,
				cfg.DB.Name,
				cfg.DB.Port,
			),
		}), &gorm.Config{})
		pool = db
		if err != nil {
			slog.Error("Unable to connect to database: %v", err)
			poolErr = err
			return
		}
	})

	return poolErr
}

func GetPool() *gorm.DB {
	return pool
}
