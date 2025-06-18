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
			slog.Error("Unable to connect to database", "error", poolErr)
			return
		}
		initDB(pool) //todo: change to versioned migration
	})

	return pool, poolErr
}

func initDB(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Room{}, &models.ChatMessage{})
	if err != nil {
		slog.Error("Error migrating database", "error", err.Error())
	}
}
