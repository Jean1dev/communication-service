package application

import (
	"encoding/json"
	"testing"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
)

func TestCreateAlert(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	alert, err := service.CreateAlert(input)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if alert.UserEmail != input.UserEmail {
		t.Errorf("Expected user_email %s, got %s", input.UserEmail, alert.UserEmail)
	}

	if alert.Type != input.Type {
		t.Errorf("Expected type %s, got %s", input.Type, alert.Type)
	}

	if !alert.Active {
		t.Error("Expected alert to be active by default")
	}
}

func TestCreateAlertWithoutEmail(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	_, err := service.CreateAlert(input)
	if err == nil {
		t.Fatal("Expected error for missing user_email, got none")
	}
}

func TestCreateAlertWithInvalidType(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      "invalid",
		Condition: condition,
	}

	_, err := service.CreateAlert(input)
	if err == nil {
		t.Fatal("Expected error for invalid type, got none")
	}
}

func TestGetAlertsByUserEmail(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	createdAlert, err := service.CreateAlert(input)
	if err != nil {
		t.Fatalf("Expected no error creating alert, got %v", err)
	}

	alerts, err := service.GetAlertsByUserEmail("test@example.com")
	if err != nil {
		t.Fatalf("Expected no error getting alerts, got %v", err)
	}

	if len(alerts) == 0 {
		t.Fatal("Expected at least one alert, got none")
	}

	found := false
	for _, alert := range alerts {
		if alert.ID == createdAlert.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find created alert in user's alerts")
	}
}

func TestToggleAlertStatus(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	createdAlert, err := service.CreateAlert(input)
	if err != nil {
		t.Fatalf("Expected no error creating alert, got %v", err)
	}

	if !createdAlert.Active {
		t.Error("Expected alert to be active initially")
	}

	toggledAlert, err := service.ToggleAlertStatus(createdAlert.ID)
	if err != nil {
		t.Fatalf("Expected no error toggling alert status, got %v", err)
	}

	if toggledAlert.Active {
		t.Error("Expected alert to be inactive after first toggle")
	}

	toggledAlert, err = service.ToggleAlertStatus(createdAlert.ID)
	if err != nil {
		t.Fatalf("Expected no error toggling alert status again, got %v", err)
	}

	if !toggledAlert.Active {
		t.Error("Expected alert to be active after second toggle")
	}
}

func TestDeleteAlert(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	createdAlert, err := service.CreateAlert(input)
	if err != nil {
		t.Fatalf("Expected no error creating alert, got %v", err)
	}

	err = service.DeleteAlert(createdAlert.ID)
	if err != nil {
		t.Fatalf("Expected no error deleting alert, got %v", err)
	}

	_, err = service.GetAlertByID(createdAlert.ID)
	if err == nil {
		t.Error("Expected error getting deleted alert, got none")
	}
}

func TestUpdateAlert(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	createdAlert, err := service.CreateAlert(input)
	if err != nil {
		t.Fatalf("Expected no error creating alert, got %v", err)
	}

	newType := dto.AlertTypeSMS
	newCondition := json.RawMessage(`{"price": {"operator": "<=", "value": 50}}`)
	updateInput := dto.UpdateAlertInput{
		Type:      &newType,
		Condition: &newCondition,
	}

	updatedAlert, err := service.UpdateAlert(createdAlert.ID, updateInput)
	if err != nil {
		t.Fatalf("Expected no error updating alert, got %v", err)
	}

	if updatedAlert.Type != newType {
		t.Errorf("Expected type %s, got %s", newType, updatedAlert.Type)
	}

	if string(updatedAlert.Condition) != string(newCondition) {
		t.Errorf("Expected condition %s, got %s", newCondition, updatedAlert.Condition)
	}
}
