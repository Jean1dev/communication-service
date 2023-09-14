package config

import (
	"communication-service/infra/database"
	"os"
)

var Db database.DefaultDatabase
var alreadyConnected = false

func connect() {
	mongouri := os.Getenv("MONGO_URI")

	if mongouri == "" {
		Db = &database.FakeRepo{}
	} else {
		Db = &database.MongoRepository{}
	}
	Db.Connect()
}

func GetDB() database.DefaultDatabase {
	if alreadyConnected == false {
		connect()
	}

	alreadyConnected = true
	return Db
}
