package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"scheduler/config"
)

var ctx context.Context
var CancelFunc context.CancelFunc
var Connection *pgxpool.Pool

func DatabaseSetUp() (err error) {
	config.SettingEnv()
	ctx, CancelFunc = context.WithTimeout(context.Background(), 100*time.Second)
	uri := fmt.Sprintf("postgres://%s/%s?user=%s&password=%s", config.PostgresHost, config.PostgresDatabase, config.PostgresUser, config.PostgresPassword)
	Connection, err = pgxpool.Connect(ctx, uri)
	return
}

func acquireConn(queryChan chan string, outputChan chan pgx.Rows, errChan chan error) {
	conn, err := Connection.Acquire(ctx)
	if err != nil {
		errChan <- err
		return
	}
	for q := range queryChan {
		r, err := conn.Query(ctx, q)
		if err != nil {
			errChan <- err
			return
		}

		for r.Next() {
			outputChan <- r
		}
		close(outputChan)
	}
}
