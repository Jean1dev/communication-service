package database

import (
	"context"
	"log"
	"os"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FakeRepo struct {
	dbMemo *MongoRepository
}

func (f *FakeRepo) Connect() {
	ctx := context.Background()
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	os.Setenv("MONGO_URI", uri)
	f.dbMemo = &MongoRepository{}
	f.dbMemo.Connect()
}

func (f *FakeRepo) Insert(data interface{}, collection string) error {
	return f.dbMemo.Insert(data, collection)
}

func (f *FakeRepo) FindAll(collection string, filter bson.D, options *options.FindOptions) (error, *mongo.Cursor) {
	return f.dbMemo.FindAll(collection, filter, options)
}

func (f *FakeRepo) UpdateOne(collection string, filter bson.D, update bson.D) error {
	return f.dbMemo.UpdateOne(collection, filter, update)
}

func (f *FakeRepo) FindOne(collection string, filter bson.D) (error, *mongo.SingleResult) {
	return f.dbMemo.FindOne(collection, filter)
}

func (f *FakeRepo) CountDocuments(collection string, filter bson.D) (int, error) {
	return f.dbMemo.CountDocuments(collection, filter)
}

func (f *FakeRepo) UpdateMany(collection string, filter bson.D, update bson.D) error {
	return f.dbMemo.UpdateMany(collection, filter, update)
}

func (f *FakeRepo) DeleteOne(collection string, filter bson.D) (int64, error) {
	return f.dbMemo.DeleteOne(collection, filter)
}
