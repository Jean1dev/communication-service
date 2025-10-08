package database

import (
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db               DefaultDatabase
	alreadyConnected = false
)

type DefaultDatabase interface {
	Connect()
	Insert(data interface{}, collection string) error
	FindAll(collection string, filter bson.D, options *options.FindOptions) (error, *mongo.Cursor)
	FindOne(collection string, filter bson.D) (error, *mongo.SingleResult)
	UpdateOne(collection string, filter bson.D, update bson.D) error
	UpdateMany(collection string, filter bson.D, update bson.D) error
	CountDocuments(collection string, filter bson.D) (int, error)
	DeleteOne(collection string, filter bson.D) (int64, error)
}

func connect() {
	mongouri := os.Getenv("MONGO_URI")

	if mongouri == "" {
		db = &FakeRepo{}
	} else {
		db = &MongoRepository{}
	}
	db.Connect()
}

func GetDB() DefaultDatabase {
	if !alreadyConnected {
		connect()
	}

	alreadyConnected = true
	return db
}
