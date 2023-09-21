package application_test

import (
	"communication-service/application"
	"testing"
)

func TestInsertNotification(t *testing.T) {
	desc := "teste"
	user := "user email"

	if err := application.InsertNewNotification(desc, user); err != nil {
		t.Errorf("InsertNewNotification throws error")
	}
}
