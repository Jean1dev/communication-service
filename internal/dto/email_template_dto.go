package dto

import (
	"errors"
	"time"
)

type EmailTemplateDto struct {
	Slug         string    `json:"slug" bson:"slug"`
	Name         string    `json:"name" bson:"name"`
	FromEmail    string    `json:"fromEmail" bson:"fromEmail"`
	HtmlTemplate string    `json:"htmlTemplate" bson:"htmlTemplate"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SendEmailV2Dto struct {
	TemplateSlug string            `json:"templateSlug"`
	To           string            `json:"to"`
	Subject      string            `json:"subject"`
	Variables    map[string]string `json:"variables"`
}

func (s *SendEmailV2Dto) Validate() error {
	if s.To == "" {
		return errors.New("to cannot be empty")
	}
	if s.Subject == "" {
		return errors.New("subject cannot be empty")
	}
	if s.TemplateSlug == "" {
		return errors.New("templateSlug cannot be empty")
	}
	return nil
}
