package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jean1dev/communication-service/api"
)

func TestRouteNotFound(t *testing.T) {
	req := httptest.NewRequest("PUT", "/notificacao", nil)
	rec := httptest.NewRecorder()

	api.NotificationHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Esperava-se status code %d, mas recebeu %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestPostNotification(t *testing.T) {
	n := api.NotificationPost{
		Desc: "teste",
		User: "teste",
	}

	data, _ := json.Marshal(n)
	req := httptest.NewRequest("POST", "/notificacao", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	api.NotificationHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperava-se status code %d, mas recebeu %d", http.StatusOK, res.StatusCode)
	}
}

func TestMarkNotificationAsRead(t *testing.T) {
	n := api.MarkNotificationAsRead{
		User: "jean@jean",
		All:  true,
	}

	data, _ := json.Marshal(n)
	req := httptest.NewRequest("POST", "/notificacao/mark-as-read", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	api.NotificationHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperava-se status code %d, mas recebeu %d", http.StatusOK, res.StatusCode)
	}
}
