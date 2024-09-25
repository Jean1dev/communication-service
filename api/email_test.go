package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jean1dev/communication-service/api"
	"github.com/Jean1dev/communication-service/internal/dto"
)

func TestEmailHandlerSuccess(t *testing.T) {
	email := dto.MailSenderInputDto{
		Recipient: "destinatario@example.com",
		Subject:   "Assunto do e-mail",
		Body:      "Conteúdo do e-mail",
	}
	emailJSON, _ := json.Marshal(email)

	req := httptest.NewRequest("POST", "/email", bytes.NewBuffer(emailJSON))
	rec := httptest.NewRecorder()

	api.EmailHandler(rec, req)

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

	api.EmailHandler(invalidRec, invalidReq)

	invalidRes := invalidRec.Result()
	defer invalidRes.Body.Close()

	if invalidRes.StatusCode != http.StatusBadRequest {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusBadRequest, invalidRes.StatusCode)
	}
}
