package common

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateDbPool(pgUrl string, l *Logger) *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, pgUrl)
	if err != nil {
		l.Error("Unable to connect to database:", err)
		os.Exit(1)
	}

	var greeting string
	err = pool.QueryRow(context.Background(), "select 'Postgres database connected.'").Scan(&greeting)
	if err != nil {
		l.Error("Fail to query db:", err)
		os.Exit(1)
	}

	l.Message(greeting)

	return pool
}

func CloseDbpool(pool *pgxpool.Pool) {
	pool.Close()
}
