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
		log.Fatalf("%s", err.Error())
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

func (m *MongoRepository) FindAll(collection string, filter bson.D, options *options.FindOptions) (error, *mongo.Cursor) {
	coll := m.db.Collection(collection)
	cursor, err := coll.Find(context.TODO(), filter, options)

	if err != nil {
		return err, nil
	}

	return nil, cursor
}

func (m *MongoRepository) UpdateOne(collection string, filter bson.D, update bson.D) error {
	coll := m.db.Collection(collection)
	result, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Print(err)
		return err
	}

	log.Printf("ModifiedCount _id: %v\n", result.ModifiedCount)
	return nil
}

func (m *MongoRepository) FindOne(collection string, filter bson.D) (error, *mongo.SingleResult) {
	coll := m.db.Collection(collection)
	doc := coll.FindOne(context.TODO(), filter)
	return nil, doc
}

func (m *MongoRepository) CountDocuments(collection string, filter bson.D) (int, error) {
	coll := m.db.Collection(collection)
	count, err := coll.CountDocuments(context.TODO(), filter)

	if err != nil {
		log.Print(err)
		return 0, err
	}

	log.Printf("CountDocuments in %s -> %d ", collection, count)
	return int(count), nil
}

func (m *MongoRepository) UpdateMany(collection string, filter bson.D, update bson.D) error {
	coll := m.db.Collection(collection)
	result, err := coll.UpdateMany(context.TODO(), filter, update)

	if err != nil {
		log.Print(err)
		return err
	}

	log.Printf("ModifiedCount : %v\n", result.ModifiedCount)
	return nil
}
