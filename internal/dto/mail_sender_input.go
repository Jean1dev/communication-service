package dto

import (
	"errors"

	"github.com/Jean1dev/communication-service/internal/templates"
)

type MailSenderInputDto struct {
	Subject         string `json:"subject"`
	Body            string `json:"message"`
	Recipient       string `json:"to"`
	TemplateCode    int    `json:"templateCode"`
	CustomBodyProps struct {
		Username    string  `json:"username"`
		Operation   string  `json:"operation"`
		Amount      float32 `json:"amount"`
		TotalAmount float32 `json:"totalAmount"`
	} `json:"customBodyProps"`
}

func (m *MailSenderInputDto) Validate() error {
	if m.Recipient == "" {
		return errors.New("recipient cannot be empty")
	}

	if m.Subject == "" {
		return errors.New("subject cannot be empty")
	}

	return nil
}

func (m *MailSenderInputDto) GetTemplate() string {
	if m.TemplateCode == 1 {
		return templates.CaixinhaTemplate(
			m.CustomBodyProps.Username,
			m.CustomBodyProps.Operation,
			m.CustomBodyProps.Amount,
			m.CustomBodyProps.TotalAmount,
		)
	}
	return templates.Default(m.Subject, m.Body)
}
