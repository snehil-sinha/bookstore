package db

import (
	"context"
	"fmt"

	"github.com/snehil-sinha/goBookStore/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// db is a MongoDB connection object
type db struct {
	Client *mongo.Client
	DB     string
}

// Client is a global connection object wrapping a Mongo Client object
// TODO: Global is bad, fix it.
var Client *db

// Connect to the db and initialise a new db instance
func New(log *common.Logger, database, uri string) (err error) {
	// Initialize the ODM object and initialise the list of Collection objects
	err = InitDbMapper(log, database, uri)
	if err != nil {
		err = fmt.Errorf("error initialising db mappers, err: %s", err.Error())
		log.Error(err.Error())
		return
	}
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(uri)

	conn, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = conn.Ping(ctx, nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	Client = &db{
		Client: conn,
		DB:     database,
	}

	return
}

func (d *db) Close(ctx context.Context, log *common.Logger) (err error) {

	if err = d.Client.Disconnect(ctx); err != nil {
		log.Error(err.Error())
	}

	return
}
