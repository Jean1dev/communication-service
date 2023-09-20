package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db *mongo.Database
}

func (m *MongoRepository) Connect() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatalf("MONGO_URI not defined")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf(err.Error())
	}

	database := client.Database("communication")
	m.db = database
	log.Print("mongo connected")
}

func (m *MongoRepository) Insert(data interface{}, collection string) error {
	coll := m.db.Collection(collection)
	result, err := coll.InsertOne(context.TODO(), data)

	if err != nil {
		log.Print(err)
		return err
	}

	log.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}

func (m *MongoRepository) FindAll(collection string, filter bson.D) (error, *mongo.Cursor) {
	coll := m.db.Collection(collection)
	cursor, err := coll.Find(context.TODO(), filter, options.Find())

	if err != nil {
		return err, nil
	}

	return nil, cursor
}
