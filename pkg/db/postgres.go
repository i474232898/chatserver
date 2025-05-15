package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool    *pgxpool.Pool
	once    sync.Once
	poolErr error
)

func Connect() (*pgxpool.Pool, error) {
	once.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			poolErr = err
			return
		}

		var result int
		err = pool.QueryRow(context.Background(), "SELECT 1").Scan(&result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection test failed: %v\n", err)
			pool.Close()
			poolErr = err
			return
		}
	})

	return pool, poolErr
}
