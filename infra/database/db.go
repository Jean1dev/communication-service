package database

import (
	"errors"
	"log"
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
	UpdateOne(collection string, filter bson.D, update bson.D) error
}

type FakeRepo struct {
}

func (f *FakeRepo) Connect() {
	log.Print("Fake repo conected")
}

func (f *FakeRepo) Insert(data interface{}, collection string) error {
	log.Printf("fake repo insert %s", data)
	return nil
}

func (f *FakeRepo) FindAll(collection string, filter bson.D, options *options.FindOptions) (error, *mongo.Cursor) {
	return errors.New("not implemented"), nil
}

func (f *FakeRepo) UpdateOne(collection string, filter bson.D, update bson.D) error {
	return errors.New("not implemented")
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
	if alreadyConnected == false {
		connect()
	}

	alreadyConnected = true
	return db
}
