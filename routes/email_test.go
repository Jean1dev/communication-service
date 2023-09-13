package routes_test

import (
	"bytes"
	"communication-service/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailHandlerSuccess(t *testing.T) {
	email := routes.Email{
		To:      "destinatario@example.com",
		Subject: "Assunto do e-mail",
		Message: "Conteúdo do e-mail",
	}
	emailJSON, _ := json.Marshal(email)

	req := httptest.NewRequest("POST", "/email", bytes.NewBuffer(emailJSON))
	rec := httptest.NewRecorder()

	routes.EmailHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 status code, but received %d", res.StatusCode)
	}
}

func TestEmailHandlerFailed(t *testing.T) {
	invalidEmailJSON := []byte(`{"to": 12, "subject": "Assunto do e-mail"}`)
	invalidReq := httptest.NewRequest("POST", "/email", bytes.NewBuffer(invalidEmailJSON))
	invalidRec := httptest.NewRecorder()

	routes.EmailHandler(invalidRec, invalidReq)

	invalidRes := invalidRec.Result()
	defer invalidRes.Body.Close()

	if invalidRes.StatusCode != http.StatusBadRequest {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusBadRequest, invalidRes.StatusCode)
	}
}
