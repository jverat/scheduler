package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"scheduler/config"
)

var client *mongo.Client
var ctx context.Context
var cancelFunc context.CancelFunc
var database *mongo.Database
var usersCollection *mongo.Collection
var mongoURI string

func DatabaseSetUp() (err error) {

	config.SettingEnv()

	mongoURI = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", config.MongoUser, config.MongoPassword, config.MongoHost, config.MongoDatabase)

	ctx, cancelFunc = context.WithTimeout(context.Background(), 100*time.Second)

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	if err != nil {
		return
	}

	database = client.Database(config.MongoDatabase)
	usersCollection = database.Collection("Users")

	return
}

func DatabaseTurnoff() (err error) {
	err = client.Disconnect(ctx)
	cancelFunc()
	return err
}
