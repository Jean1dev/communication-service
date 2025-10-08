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
