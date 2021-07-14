package db

import (
	"context"
	"fmt"
	"time"

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
	if err != nil {
		return err
	}
	return
}

func prepareStatements() (err error) {
	acquire, err := Connection.Acquire(ctx)
	if err != nil {
		return err
	}
	_, err = acquire.Conn().Prepare(ctx, "Read", `SELECT * FROM user !1`)
	if err != nil {
		return err
	}
	return
}
