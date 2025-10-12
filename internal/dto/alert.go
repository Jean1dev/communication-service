package dto

import (
	"encoding/json"
	"time"
)

type AlertType string

const (
	AlertTypeSMS      AlertType = "sms"
	AlertTypeEmail    AlertType = "email"
	AlertTypeTelegram AlertType = "telegram"
)

type Alert struct {
	ID        string          `json:"id" bson:"_id,omitempty"`
	UserEmail string          `json:"user_email" bson:"user_email"`
	Type      AlertType       `json:"type" bson:"type"`
	Condition json.RawMessage `json:"condition" bson:"condition"`
	Active    bool            `json:"active" bson:"active"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
}

type CreateAlertInput struct {
	UserEmail string          `json:"user_email"`
	Type      AlertType       `json:"type"`
	Condition json.RawMessage `json:"condition"`
}

type AlertResponse struct {
	ID        string          `json:"id"`
	UserEmail string          `json:"user_email"`
	Type      AlertType       `json:"type"`
	Condition json.RawMessage `json:"condition"`
	Active    bool            `json:"active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type GroupedAlertsResponse struct {
	UserEmail string          `json:"user_email"`
	Alerts    []AlertResponse `json:"alerts"`
}

func (a *Alert) ToResponse() AlertResponse {
	return AlertResponse{
		ID:        a.ID,
		UserEmail: a.UserEmail,
		Type:      a.Type,
		Condition: a.Condition,
		Active:    a.Active,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
