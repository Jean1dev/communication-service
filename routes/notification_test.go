package routes_test

import (
	"bytes"
	"communication-service/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteNotFound(t *testing.T) {
	req := httptest.NewRequest("PUT", "/notificacao", nil)
	rec := httptest.NewRecorder()

	routes.NotificationHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Esperava-se status code %d, mas recebeu %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestPostNotification(t *testing.T) {
	n := routes.NotificationPost{
		Desc: "teste",
		User: "teste",
	}

	data, _ := json.Marshal(n)
	req := httptest.NewRequest("POST", "/notificacao", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	routes.NotificationHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperava-se status code %d, mas recebeu %d", http.StatusOK, res.StatusCode)
	}
}

func TestMarkNotificationAsRead(t *testing.T) {
	n := routes.MarkNotificationAsRead{
		User: "jean@jean",
		All:  true,
	}

	data, _ := json.Marshal(n)
	req := httptest.NewRequest("POST", "/notificacao/mark-as-read", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	routes.NotificationHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperava-se status code %d, mas recebeu %d", http.StatusOK, res.StatusCode)
	}
}
