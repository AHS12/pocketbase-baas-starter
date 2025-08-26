package jobs

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"os"
	"path/filepath"

	"ims-pocketbase-baas-starter/pkg/cronutils"
	"ims-pocketbase-baas-starter/pkg/jobutils"
	log "ims-pocketbase-baas-starter/pkg/logger"
	"ims-pocketbase-baas-starter/pkg/metrics"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

// EmailJobHandler handles email job processing
type EmailJobHandler struct {
	app *pocketbase.PocketBase
}

// NewEmailJobHandler creates a new email job handler
func NewEmailJobHandler(app *pocketbase.PocketBase) *EmailJobHandler {
	return &EmailJobHandler{
		app: app,
	}
}

// Handle processes an email job using typed payload structures (with metrics instrumentation)
func (h *EmailJobHandler) Handle(ctx *cronutils.CronExecutionContext, job *jobutils.JobData) error {
	ctx.LogStart(fmt.Sprintf("Processing email job: %s", job.ID))

	metricsProvider := metrics.GetInstance()

	// Instrument the job handler execution with metrics collection
	return metrics.InstrumentJobHandler(metricsProvider, "email_job", func() error {
		emailPayload, err := jobutils.ParseEmailJobPayload(job)
		if err != nil {
			return fmt.Errorf("failed to parse email job payload: %w", err)
		}

		if err := h.validateEmailPayload(emailPayload); err != nil {
			return fmt.Errorf("invalid email job payload: %w", err)
		}

		htmlContent, textContent, err := h.processEmailTemplates(emailPayload)
		if err != nil {
			return fmt.Errorf("failed to process email templates: %w", err)
		}

		if err := h.sendEmail(emailPayload, htmlContent, textContent); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		ctx.LogEnd("Email job processed successfully")
		return nil
	})
}

// GetJobType returns the job type this handler processes
func (h *EmailJobHandler) GetJobType() string {
	return jobutils.JobTypeEmail
}

// validateEmailPayload validates the typed email job payload (additional handler-specific validation)
func (h *EmailJobHandler) validateEmailPayload(payload *jobutils.EmailJobPayload) error {
	if payload.Type != jobutils.JobTypeEmail {
		return fmt.Errorf("invalid job type: expected %s, got %s", jobutils.JobTypeEmail, payload.Type)
	}

	// Additional validation can be added here for handler-specific requirements

	return nil
}

// processEmailTemplates processes both HTML and text email templates with variables
func (h *EmailJobHandler) processEmailTemplates(payload *jobutils.EmailJobPayload) (string, string, error) {
	if payload.Data.Template == "" {
		log.Warn("No template specified, using empty content")
		return "", "", nil
	}

	htmlContent, err := h.processSingleTemplate(payload, ".html")
	if err != nil {
		log.Warn("Failed to process HTML template", "error", err)
	}

	textContent, err := h.processSingleTemplate(payload, ".txt")
	if err != nil {
		log.Warn("Failed to process text template", "error", err)
	}

	return htmlContent, textContent, nil
}

// processSingleTemplate processes a single email template with variables
func (h *EmailJobHandler) processSingleTemplate(payload *jobutils.EmailJobPayload, extension string) (string, error) {
	templatePath := filepath.Join("templates", "emails", payload.Data.Template+extension)

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, payload.Data.Variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// sendEmail sends the email using PocketBase mailer
func (h *EmailJobHandler) sendEmail(payload *jobutils.EmailJobPayload, htmlContent, textContent string) error {
	settings := h.app.Settings()

	// Use the configured sender name and address from admin UI
	fromEmail := settings.Meta.SenderAddress
	fromName := settings.Meta.SenderName

	// Fallback to environment variables if admin UI settings are empty
	if fromEmail == "" {
		fromEmail = os.Getenv("SMTP_FROM_EMAIL")
		if fromEmail == "" {
			fromEmail = "noreply@ims-app.local"
		}
	}

	if fromName == "" {
		fromName = os.Getenv("SMTP_FROM_NAME")
		if fromName == "" {
			fromName = "IMS PocketBase App"
		}
	}

	message := &mailer.Message{
		From:    mail.Address{Name: fromName, Address: fromEmail},
		To:      []mail.Address{{Address: payload.Data.To}},
		Subject: payload.Data.Subject,
	}

	if htmlContent != "" {
		message.HTML = htmlContent
	}

	if textContent != "" {
		message.Text = textContent
	} else if htmlContent != "" {
		// If only HTML content is available, use it as text as fallback
		message.Text = htmlContent
	}

	if err := h.app.NewMailClient().Send(message); err != nil {
		log.Error("Failed to send email",
			"to", payload.Data.To,
			"subject", payload.Data.Subject,
			"error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info("Email sent successfully",
		"to", payload.Data.To,
		"subject", payload.Data.Subject)

	return nil
}
