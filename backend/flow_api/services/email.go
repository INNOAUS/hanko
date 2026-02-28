package services

import (
	"fmt"
	"strings"

	"github.com/teamhanko/hanko/backend/v2/config"
	"github.com/teamhanko/hanko/backend/v2/mail"
	"gopkg.in/gomail.v2"
)

// injectMailBranding adds mail_template.yaml branding for the given lang to the template data.
// When set, layout uses these; otherwise falls back to locale translations.
// ServiceName is also set to product_name when present, so body text (e.g. "您的 {{ .ServiceName }} 帐户...") shows the same product name.
func injectMailBranding(data map[string]interface{}, lang string, tmpl *config.MailTemplateConfig) {
	if data == nil || tmpl == nil {
		return
	}
	b := tmpl.BrandingForLang(lang)
	if strings.TrimSpace(b.ProductName) != "" {
		data["MailProductName"] = b.ProductName
		data["ServiceName"] = b.ProductName
	}
	if strings.TrimSpace(b.FooterSentBy) != "" {
		data["MailFooterSentBy"] = b.FooterSentBy
	}
	if strings.TrimSpace(b.Copyright) != "" {
		data["MailCopyright"] = b.Copyright
	}
}

type Email struct {
	renderer     *mail.Renderer
	mailer       mail.Mailer
	cfg          config.Config
	mailTemplate *config.MailTemplateConfig
}

// NewEmailService creates the email service. mailTemplate may be nil (then locale files are used for layout and body).
func NewEmailService(cfg config.Config, mailTemplate *config.MailTemplateConfig) (*Email, error) {
	store := mail.NewMailTemplateStoreFromConfig(mailTemplate)
	renderer, err := mail.NewRenderer(store)
	if err != nil {
		return nil, err
	}
	mailer, err := mail.NewMailer(cfg.EmailDelivery.SMTP)
	if err != nil {
		panic(fmt.Errorf("failed to create mailer: %w", err))
	}

	return &Email{
		renderer:     renderer,
		mailer:      mailer,
		cfg:         cfg,
		mailTemplate: mailTemplate,
	}, nil
}

// SendEmail sends an email to the emailAddress with the given subject and body.
func (s *Email) SendEmail(emailAddress, subject, body, htmlBody string) error {
	message := gomail.NewMessage()
	message.SetAddressHeader("To", emailAddress, "")
	message.SetAddressHeader("From", s.cfg.EmailDelivery.FromAddress, s.cfg.EmailDelivery.FromName)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)
	message.AddAlternative("text/html", htmlBody)

	if err := s.mailer.Send(message); err != nil {
		return err
	}

	return nil
}

// RenderSubject renders a subject with the given template. Must be "subject_[template_name]".
// Injects mail_template branding so subject strings using {{ .ServiceName }} use the configured product name.
func (s *Email) RenderSubject(lang, template string, data map[string]interface{}) string {
	if data != nil {
		injectMailBranding(data, lang, s.mailTemplate)
	}
	return s.renderer.Translate(lang, fmt.Sprintf("subject_%s", template), data)
}

// RenderBodyPlain renders the body with the given template. The template name must be the name of the template without the
// content type and the file ending. E.g. when the file is created as "email_verification_text.tmpl" then the template
// name is just "email_verification"
func (s *Email) RenderBodyPlain(lang, template string, data map[string]interface{}) (string, error) {
	injectMailBranding(data, lang, s.mailTemplate)
	return s.renderer.RenderPlain(template, lang, data)
}

func (s *Email) RenderBodyHTML(lang, template string, data map[string]interface{}) (string, error) {
	injectMailBranding(data, lang, s.mailTemplate)
	return s.renderer.RenderHTML(template, lang, data)
}
