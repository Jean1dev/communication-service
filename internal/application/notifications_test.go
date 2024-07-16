package application_test

import (
	"testing"

	"github.com/Jean1dev/communication-service/internal/application"
)

func TestInsertNotification(t *testing.T) {
	desc := "teste"
	user := "user email"

	if err := application.InsertNewNotification(desc, user); err != nil {
		t.Errorf("InsertNewNotification throws error")
	}
}
