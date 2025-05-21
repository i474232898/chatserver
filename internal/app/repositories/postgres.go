package repositories

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/i474232898/chatserver/configs"
	"github.com/i474232898/chatserver/internal/app/repositories/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pool    *gorm.DB
	once    sync.Once
	poolErr error
)

func GetPool(cfg *configs.AppConfigs) (*gorm.DB, error) {
	once.Do(func() {
		pool, poolErr = gorm.Open(postgres.New(postgres.Config{
			DSN: fmt.Sprintf(
				"user=%s host=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
				cfg.DB.User,
				cfg.DB.Host,
				cfg.DB.Password,
				cfg.DB.Name,
				cfg.DB.Port,
			),
		}), &gorm.Config{})
		if poolErr != nil {
			slog.Error("Unable to connect to database: %v", poolErr)
			return
		}
		initDB(pool) //todo: change to versioned migration 
	})

	return pool, poolErr
}

func initDB(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
