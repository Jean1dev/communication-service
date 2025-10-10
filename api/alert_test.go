package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
)

func setupTestDB() {
	db := &database.FakeRepo{}
	db.Connect()
	alertService = nil
}

func TestAlertHandlerCreate(t *testing.T) {
	setupTestDB()

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AlertHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response dto.AlertResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.UserEmail != input.UserEmail {
		t.Errorf("Expected user_email %s, got %s", input.UserEmail, response.UserEmail)
	}

	if response.Type != input.Type {
		t.Errorf("Expected type %s, got %s", input.Type, response.Type)
	}

	if !response.Active {
		t.Error("Expected alert to be active")
	}
}

func TestAlertHandlerCreateWithoutEmail(t *testing.T) {
	setupTestDB()

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AlertHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAlertHandlerGetByUserEmail(t *testing.T) {
	setupTestDB()

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	body, _ := json.Marshal(input)
	reqCreate := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	AlertHandler(wCreate, reqCreate)

	req := httptest.NewRequest(http.MethodGet, "/alerts?user_email=test@example.com", nil)
	w := httptest.NewRecorder()

	AlertHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []dto.AlertResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) == 0 {
		t.Error("Expected at least one alert")
	}
}

func TestAlertHandlerGetWithoutUserEmail(t *testing.T) {
	setupTestDB()

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	w := httptest.NewRecorder()

	AlertHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAlertToggleStatusHandler(t *testing.T) {
	setupTestDB()

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	body, _ := json.Marshal(input)
	reqCreate := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	AlertHandler(wCreate, reqCreate)

	var createResponse dto.AlertResponse
	json.NewDecoder(wCreate.Body).Decode(&createResponse)

	if !createResponse.Active {
		t.Error("Expected alert to be active initially")
	}

	req := httptest.NewRequest(http.MethodPost, "/alerts/toggle/"+createResponse.ID, nil)
	w := httptest.NewRecorder()

	AlertToggleStatusHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var toggleResponse dto.AlertResponse
	json.NewDecoder(w.Body).Decode(&toggleResponse)

	if toggleResponse.Active {
		t.Error("Expected alert to be inactive after first toggle")
	}

	req2 := httptest.NewRequest(http.MethodPost, "/alerts/toggle/"+createResponse.ID, nil)
	w2 := httptest.NewRecorder()

	AlertToggleStatusHandler(w2, req2)

	var toggleResponse2 dto.AlertResponse
	json.NewDecoder(w2.Body).Decode(&toggleResponse2)

	if !toggleResponse2.Active {
		t.Error("Expected alert to be active after second toggle")
	}
}

func TestAlertHandlerDelete(t *testing.T) {
	setupTestDB()

	condition := json.RawMessage(`{"price": {"operator": ">=", "value": 100}}`)
	input := dto.CreateAlertInput{
		UserEmail: "test@example.com",
		Type:      dto.AlertTypeEmail,
		Condition: condition,
	}

	body, _ := json.Marshal(input)
	reqCreate := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	AlertHandler(wCreate, reqCreate)

	var createResponse dto.AlertResponse
	json.NewDecoder(wCreate.Body).Decode(&createResponse)

	req := httptest.NewRequest(http.MethodDelete, "/alerts/"+createResponse.ID, nil)
	w := httptest.NewRecorder()

	AlertHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestAlertHandlerMethodNotAllowed(t *testing.T) {
	setupTestDB()

	req := httptest.NewRequest(http.MethodPatch, "/alerts", nil)
	w := httptest.NewRecorder()

	AlertHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestAlertsGroupedHandler(t *testing.T) {
	setupTestDB()

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

	body1, _ := json.Marshal(input1)
	reqCreate1 := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body1))
	reqCreate1.Header.Set("Content-Type", "application/json")
	wCreate1 := httptest.NewRecorder()
	AlertHandler(wCreate1, reqCreate1)

	body2, _ := json.Marshal(input2)
	reqCreate2 := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body2))
	reqCreate2.Header.Set("Content-Type", "application/json")
	wCreate2 := httptest.NewRecorder()
	AlertHandler(wCreate2, reqCreate2)

	body3, _ := json.Marshal(input3)
	reqCreate3 := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body3))
	reqCreate3.Header.Set("Content-Type", "application/json")
	wCreate3 := httptest.NewRecorder()
	AlertHandler(wCreate3, reqCreate3)

	req := httptest.NewRequest(http.MethodGet, "/alerts/grouped", nil)
	w := httptest.NewRecorder()

	AlertsGroupedHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []dto.GroupedAlertsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) < 2 {
		t.Errorf("Expected at least 2 groups, got %d", len(response))
	}

	user1Count := 0
	user2Count := 0

	for _, group := range response {
		if group.UserEmail == "user1@example.com" {
			user1Count = len(group.Alerts)
		}
		if group.UserEmail == "user2@example.com" {
			user2Count = len(group.Alerts)
		}
	}

	if user1Count < 2 {
		t.Errorf("Expected at least 2 alerts for user1@example.com, got %d", user1Count)
	}

	if user2Count < 1 {
		t.Errorf("Expected at least 1 alert for user2@example.com, got %d", user2Count)
	}
}

func TestAlertsGroupedHandlerOnlyActive(t *testing.T) {
	setupTestDB()

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

	body1, _ := json.Marshal(input1)
	reqCreate1 := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body1))
	reqCreate1.Header.Set("Content-Type", "application/json")
	wCreate1 := httptest.NewRecorder()
	AlertHandler(wCreate1, reqCreate1)

	var alert1Response dto.AlertResponse
	json.NewDecoder(wCreate1.Body).Decode(&alert1Response)

	body2, _ := json.Marshal(input2)
	reqCreate2 := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(body2))
	reqCreate2.Header.Set("Content-Type", "application/json")
	wCreate2 := httptest.NewRecorder()
	AlertHandler(wCreate2, reqCreate2)

	var alert2Response dto.AlertResponse
	json.NewDecoder(wCreate2.Body).Decode(&alert2Response)

	reqToggle := httptest.NewRequest(http.MethodPost, "/alerts/toggle/"+alert2Response.ID, nil)
	wToggle := httptest.NewRecorder()
	AlertToggleStatusHandler(wToggle, reqToggle)

	req := httptest.NewRequest(http.MethodGet, "/alerts/grouped", nil)
	w := httptest.NewRecorder()

	AlertsGroupedHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []dto.GroupedAlertsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	user3Count := 0
	for _, group := range response {
		if group.UserEmail == "user3@example.com" {
			user3Count = len(group.Alerts)
			for _, alert := range group.Alerts {
				if !alert.Active {
					t.Error("Expected only active alerts, found inactive alert")
				}
				if alert.ID == alert2Response.ID {
					t.Error("Expected inactive alert to be excluded from results")
				}
			}
		}
	}

	if user3Count != 1 {
		t.Errorf("Expected 1 active alert for user3@example.com, got %d", user3Count)
	}
}
