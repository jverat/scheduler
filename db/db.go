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

func AcquireConn(queryChan chan string, outputChan chan pgx.Rows, errChan chan error) {
	conn, err := Connection.Acquire(ctx)
	if err != nil {
		errChan <- err
		if conn != nil {
			conn.Release()
		}
		return
	}

	for q := range queryChan {
		r, err := conn.Query(ctx, q)
		if err != nil {
			errChan <- err
		} else {
			outputChan <- r
		}
	}
	close(outputChan)
	conn.Release()
}
