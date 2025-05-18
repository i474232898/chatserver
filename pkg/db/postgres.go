package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/i474232898/chatserver/configs"
)

var (
	pool    *pgxpool.Pool
	once    sync.Once
	poolErr error
)

func Connect(cfg *configs.AppConfigs) error {
	once.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			poolErr = err
			return
		}
	})

	return poolErr
}

func GetPool() *pgxpool.Pool {
	return pool
}
