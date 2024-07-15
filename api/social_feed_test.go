package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Jean1dev/communication-service/api"
)

func TestNaoDeveCriarUmNovoPos(t *testing.T) {
	body := api.SocialPost{
		Message: "my first post",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/social-feed", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	api.SocialFeedHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Esperava http.StatusBadRequest obtive outra coisa")
	}
}
