package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	once sync.Once
	conn *pgx.Conn
	err  error
)

func Connect() (*pgx.Conn, error) {
	once.Do(func() {
		conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		}

		var result int
		err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Connection test failed: %v\n", err)
			os.Exit(1)
		}
	})

	return conn, err
}
