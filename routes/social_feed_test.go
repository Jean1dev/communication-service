package routes_test

import (
	"bytes"
	"communication-service/routes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNaoDeveCriarUmNovoPos(t *testing.T) {
	body := routes.SocialPost{
		Message: "my first post",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/social-feed", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	routes.SocialFeedHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Esperava http.StatusBadRequest obtive outra coisa")
	}
}
