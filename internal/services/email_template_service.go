package services

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"github.com/Jean1dev/communication-service/internal/dto"
	"github.com/Jean1dev/communication-service/internal/infra/database"
	"go.mongodb.org/mongo-driver/bson"
)

const emailTemplatesCollection = "email_templates_v2"

func CreateEmailTemplate(t dto.EmailTemplateDto) error {
	if t.Slug == "" {
		return errors.New("slug cannot be empty")
	}
	if t.HtmlTemplate == "" {
		return errors.New("htmlTemplate cannot be empty")
	}

	db := database.GetDB()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return db.Insert(t, emailTemplatesCollection)
}

func GetEmailTemplate(slug string) (*dto.EmailTemplateDto, error) {
	db := database.GetDB()
	filter := bson.D{{Key: "slug", Value: slug}}
	err, result := db.FindOne(emailTemplatesCollection, filter)
	if err != nil {
		return nil, err
	}

	var tmpl dto.EmailTemplateDto
	if err := result.Decode(&tmpl); err != nil {
		return nil, errors.New("template not found")
	}
	return &tmpl, nil
}

func ListEmailTemplates() ([]dto.EmailTemplateDto, error) {
	db := database.GetDB()
	err, cursor := db.FindAll(emailTemplatesCollection, bson.D{}, nil)
	if err != nil {
		return nil, err
	}

	var templates []dto.EmailTemplateDto
	if err := cursor.All(context.TODO(), &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func UpdateEmailTemplate(slug string, updated dto.EmailTemplateDto) error {
	db := database.GetDB()
	filter := bson.D{{Key: "slug", Value: slug}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: updated.Name},
		{Key: "fromEmail", Value: updated.FromEmail},
		{Key: "htmlTemplate", Value: updated.HtmlTemplate},
		{Key: "updatedAt", Value: time.Now()},
	}}}
	return db.UpdateOne(emailTemplatesCollection, filter, update)
}

func DeleteEmailTemplate(slug string) error {
	db := database.GetDB()
	filter := bson.D{{Key: "slug", Value: slug}}
	_, err := db.DeleteOne(emailTemplatesCollection, filter)
	return err
}

func SendEmailV2(input dto.SendEmailV2Dto) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tmpl, err := GetEmailTemplate(input.TemplateSlug)
	if err != nil {
		return err
	}

	html, err := renderEmailTemplate(tmpl.HtmlTemplate, input.Variables)
	if err != nil {
		return err
	}

	fromEmail := tmpl.FromEmail
	if fromEmail == "" {
		fromEmail = "info@jeanconsultoria.com"
	}

	AsyncSendRaw(input.Subject, input.To, html, fromEmail, "")
	return nil
}

func renderEmailTemplate(htmlTemplate string, variables map[string]string) (string, error) {
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	data := make(map[string]interface{}, len(variables))
	for k, v := range variables {
		data[k] = v
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}
