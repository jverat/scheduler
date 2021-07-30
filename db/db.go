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
	ctx, CancelFunc = context.WithTimeout(context.Background(), 1000*time.Second)
	uri := fmt.Sprintf("postgres://%s/%s?user=%s&password=%s", config.PostgresHost, config.PostgresDatabase, config.PostgresUser, config.PostgresPassword)
	Connection, err = pgxpool.Connect(ctx, uri)
	return
}

func AcquireConn(queryChan chan string, outputChan chan pgx.Rows, errChan chan error) {

	for q := range queryChan {

		conn, err := Connection.Acquire(ctx)
		if err != nil {
			errChan <- err
			if conn != nil {
				conn.Release()
			}
			return
		}

		r, err := conn.Query(ctx, q)
		conn.Release()
		if err != nil {
			errChan <- err
			return
		} else {
			outputChan <- r
		}

	}
	close(outputChan)
}
