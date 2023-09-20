package database

import (
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultDatabase interface {
	Connect()
	Insert(data interface{}, collection string) error
	FindAll(collection string, filter bson.D) (error, *mongo.Cursor)
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

func (f *FakeRepo) FindAll(collection string, filter bson.D) (error, *mongo.Cursor) {
	return errors.New("not implemented"), nil
}
