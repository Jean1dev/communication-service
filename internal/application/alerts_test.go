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

func TestGetAllAlertsGroupedByEmail(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition1 := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input1 := dto.CreateAlertInput{
		UserEmail: "user1@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition1,
	}

	condition2 := json.RawMessage(`{"price": {"operator": "<=", "value": 50}}`)
	input2 := dto.CreateAlertInput{
		UserEmail: "user1@example.com",
		Type:      dto.AlertTypeSMS,
		Condition: condition2,
	}

	condition3 := json.RawMessage(`{"price": {"operator": ">=", "value": 200}}`)
	input3 := dto.CreateAlertInput{
		UserEmail: "user2@example.com",
		Type:      dto.AlertTypeTelegram,
		Condition: condition3,
	}

	_, err := service.CreateAlert(input1)
	if err != nil {
		t.Fatalf("Expected no error creating alert 1, got %v", err)
	}

	_, err = service.CreateAlert(input2)
	if err != nil {
		t.Fatalf("Expected no error creating alert 2, got %v", err)
	}

	_, err = service.CreateAlert(input3)
	if err != nil {
		t.Fatalf("Expected no error creating alert 3, got %v", err)
	}

	groupedAlerts, err := service.GetAllAlertsGroupedByEmail()
	if err != nil {
		t.Fatalf("Expected no error getting grouped alerts, got %v", err)
	}

	if len(groupedAlerts) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groupedAlerts))
	}

	user1Count := 0
	user2Count := 0

	for _, group := range groupedAlerts {
		if group.UserEmail == "user1@example.com" {
			user1Count = len(group.Alerts)
		}
		if group.UserEmail == "user2@example.com" {
			user2Count = len(group.Alerts)
		}
	}

	if user1Count != 2 {
		t.Errorf("Expected 2 alerts for user1@example.com, got %d", user1Count)
	}

	if user2Count != 1 {
		t.Errorf("Expected 1 alert for user2@example.com, got %d", user2Count)
	}
}

func TestGetAllAlertsGroupedByEmailOnlyActive(t *testing.T) {
	db := &database.FakeRepo{}
	db.Connect()
	service := NewAlertService(db)

	condition1 := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input1 := dto.CreateAlertInput{
		UserEmail: "user3@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition1,
	}

	condition2 := json.RawMessage(`{"price": {"operator": "<=", "value": 50}}`)
	input2 := dto.CreateAlertInput{
		UserEmail: "user3@example.com",
		Type:      dto.AlertTypeSMS,
		Condition: condition2,
	}

	alert1, err := service.CreateAlert(input1)
	if err != nil {
		t.Fatalf("Expected no error creating alert 1, got %v", err)
	}

	alert2, err := service.CreateAlert(input2)
	if err != nil {
		t.Fatalf("Expected no error creating alert 2, got %v", err)
	}

	_, err = service.ToggleAlertStatus(alert2.ID)
	if err != nil {
		t.Fatalf("Expected no error toggling alert status, got %v", err)
	}

	groupedAlerts, err := service.GetAllAlertsGroupedByEmail()
	if err != nil {
		t.Fatalf("Expected no error getting grouped alerts, got %v", err)
	}

	user3Count := 0
	for _, group := range groupedAlerts {
		if group.UserEmail == "user3@example.com" {
			user3Count = len(group.Alerts)
			for _, alert := range group.Alerts {
				if !alert.Active {
					t.Error("Expected only active alerts, found inactive alert")
				}
				if alert.ID == alert2.ID {
					t.Error("Expected inactive alert to be excluded from results")
				}
				if alert.ID == alert1.ID && !alert.Active {
					t.Error("Expected active alert to be included and active")
				}
			}
		}
	}

	if user3Count != 1 {
		t.Errorf("Expected 1 active alert for user3@example.com, got %d", user3Count)
	}
}
