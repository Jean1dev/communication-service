package dto

import (
	"errors"

	"github.com/Jean1dev/communication-service/internal/templates"
)

type CustomBodyPropsDto struct {
	Username    string  `json:"username"`
	Operation   string  `json:"operation"`
	Amount      float32 `json:"amount"`
	TotalAmount float32 `json:"totalAmount"`

	ChavePix      string `json:"chavePix"`
	PixCopiaECola string `json:"pixCopiaECola"`
	PixQRCode     string `json:"pixQRCode"`
	LinkPagamento string `json:"linkPagamento"`

	MainMessage string `json:"mainMessage"`
	ContactLink string `json:"contactLink"`
	PrivacyLink string `json:"privacyLink"`
	TermsLink   string `json:"termsLink"`
}

type MailSenderInputDto struct {
	Subject         string             `json:"subject"`
	Body            string             `json:"message"`
	Recipient       string             `json:"to"`
	TemplateCode    int                `json:"templateCode"`
	AttachmentLink  string             `json:"attachmentLink"`
	CustomBodyProps CustomBodyPropsDto `json:"customBodyProps"`
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

	if m.TemplateCode == 2 {
		return templates.PaymentTemplate(
			m.CustomBodyProps.LinkPagamento,
			m.CustomBodyProps.ChavePix,
			m.CustomBodyProps.PixCopiaECola,
			m.CustomBodyProps.Username,
			m.Body,
		)
	}

	if m.TemplateCode == 3 {
		return templates.MeConecteiTemplate(
			m.CustomBodyProps.Username,
			m.Recipient,
			m.CustomBodyProps.MainMessage,
			m.CustomBodyProps.ContactLink,
			m.CustomBodyProps.PrivacyLink,
			m.CustomBodyProps.TermsLink,
		)
	}
	return templates.Default(m.Subject, m.Body)
}
