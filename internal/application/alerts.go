package application

import (
	"errors"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
)

type AlertService struct {
	alertRepo *database.AlertRepository
}

func NewAlertService(db database.DefaultDatabase) *AlertService {
	return &AlertService{
		alertRepo: database.NewAlertRepository(db),
	}
}

func (s *AlertService) CreateAlert(input dto.CreateAlertInput) (*dto.Alert, error) {
	if input.UserEmail == "" {
		return nil, errors.New("user_email is required")
	}

	if input.Type != dto.AlertTypeSMS && input.Type != dto.AlertTypeEmail && input.Type != dto.AlertTypeTelegram {
		return nil, errors.New("invalid alert type, must be sms, email or telegram")
	}

	if len(input.Condition) == 0 {
		return nil, errors.New("condition is required")
	}

	alert := &dto.Alert{
		UserEmail: input.UserEmail,
		Type:      input.Type,
		Condition: input.Condition,
	}

	createdAlert, err := s.alertRepo.Create(alert)
	if err != nil {
		return nil, err
	}

	return createdAlert, nil
}

func (s *AlertService) GetAlertByID(id string) (*dto.Alert, error) {
	if id == "" {
		return nil, errors.New("alert id is required")
	}

	alert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return alert, nil
}

func (s *AlertService) GetAlertsByUserEmail(userEmail string) ([]dto.Alert, error) {
	if userEmail == "" {
		return nil, errors.New("user_email is required")
	}

	alerts, err := s.alertRepo.FindByUserEmail(userEmail)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

func (s *AlertService) UpdateAlert(id string, input dto.UpdateAlertInput) (*dto.Alert, error) {
	if id == "" {
		return nil, errors.New("alert id is required")
	}

	alert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Type != nil {
		if *input.Type != dto.AlertTypeSMS && *input.Type != dto.AlertTypeEmail && *input.Type != dto.AlertTypeTelegram {
			return nil, errors.New("invalid alert type, must be sms, email or telegram")
		}
		alert.Type = *input.Type
	}

	if input.Condition != nil {
		if len(*input.Condition) == 0 {
			return nil, errors.New("condition cannot be empty")
		}
		alert.Condition = *input.Condition
	}

	err = s.alertRepo.Update(id, alert)
	if err != nil {
		return nil, err
	}

	updatedAlert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedAlert, nil
}

func (s *AlertService) DeleteAlert(id string) error {
	if id == "" {
		return errors.New("alert id is required")
	}

	_, err := s.alertRepo.FindByID(id)
	if err != nil {
		return err
	}

	err = s.alertRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AlertService) ToggleAlertStatus(id string) (*dto.Alert, error) {
	if id == "" {
		return nil, errors.New("alert id is required")
	}

	alert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	err = s.alertRepo.SetActive(id, !alert.Active)
	if err != nil {
		return nil, err
	}

	updatedAlert, err := s.alertRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return updatedAlert, nil
}
