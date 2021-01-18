package web

import (
	"context"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbDatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Passowrd string
	Database string
}


func NewMongoDatabase(config MongoDbDatabaseConfig) (*mongo.Database, error) {
	if config.Host == "" {
		return nil, nil
	}
	dblink := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		config.User,
		config.Passowrd,
		config.Host,
		config.Port,
	)
	client, err := mongo.NewClient(options.Client().ApplyURI(dblink))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return client.Database(config.Database), nil
}

func NewMongoDatabaseFromConfig(configName string) (*mongo.Database, error) {
	dbhost := globalBasic.Config.GetString(configName + "host")
	dbport := globalBasic.Config.GetString(configName + "port")
	dbuser := globalBasic.Config.GetString(configName + "user")
	dbpassword := globalBasic.Config.GetString(configName + "password")
	dbdatabase := globalBasic.Config.GetString(configName + "database")

	config := MongoDbDatabaseConfig{}
	config.Host = dbhost
	config.Port, _ = strconv.Atoi(dbport)
	config.User = dbuser
	config.Passowrd = dbpassword
	config.Database = dbdatabase

	return NewMongoDatabase(config)
}
