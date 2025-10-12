package database

import (
	"context"
	"errors"
	"time"

	"github.com/Jean1dev/communication-service/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const AlertsCollection = "user_alerts"

type AlertRepository struct {
	db DefaultDatabase
}

func NewAlertRepository(db DefaultDatabase) *AlertRepository {
	return &AlertRepository{
		db: db,
	}
}

func (r *AlertRepository) Create(alert *dto.Alert) (*dto.Alert, error) {
	alert.ID = primitive.NewObjectID().Hex()
	alert.CreatedAt = time.Now()
	alert.UpdatedAt = time.Now()
	alert.Active = true

	err := r.db.Insert(alert, AlertsCollection)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (r *AlertRepository) FindByID(id string) (*dto.Alert, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	_, result := r.db.FindOne(AlertsCollection, filter)

	var alert dto.Alert
	err := result.Decode(&alert)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("alert not found")
		}
		return nil, err
	}

	return &alert, nil
}

func (r *AlertRepository) FindByUserEmail(userEmail string) ([]dto.Alert, error) {
	filter := bson.D{{Key: "user_email", Value: userEmail}}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	_, cursor := r.db.FindAll(AlertsCollection, filter, opts)
	if cursor == nil {
		return []dto.Alert{}, nil
	}
	defer cursor.Close(context.TODO())

	var alerts []dto.Alert
	if err := cursor.All(context.TODO(), &alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (r *AlertRepository) FindAll() ([]dto.Alert, error) {
	filter := bson.D{{Key: "active", Value: true}}
	opts := options.Find().SetSort(bson.D{{Key: "user_email", Value: 1}, {Key: "created_at", Value: -1}})

	_, cursor := r.db.FindAll(AlertsCollection, filter, opts)
	if cursor == nil {
		return []dto.Alert{}, nil
	}
	defer cursor.Close(context.TODO())

	var alerts []dto.Alert
	if err := cursor.All(context.TODO(), &alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

func (r *AlertRepository) Delete(id string) error {
	filter := bson.D{{Key: "_id", Value: id}}

	deletedCount, err := r.db.DeleteOne(AlertsCollection, filter)
	if err != nil {
		return err
	}

	if deletedCount == 0 {
		return errors.New("alert not found")
	}

	return nil
}

func (r *AlertRepository) SetActive(id string, active bool) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "active", Value: active},
		{Key: "updated_at", Value: time.Now()},
	}}}

	err := r.db.UpdateOne(AlertsCollection, filter, update)
	if err != nil {
		return err
	}

	return nil
}
