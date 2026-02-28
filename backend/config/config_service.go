package config

import (
	"errors"
	"strings"
)

// MailBranding holds optional overrides for mail layout text per locale.
// Used by mail_template.yaml (keyed by language code).
type MailBranding struct {
	// Product name shown in the email header.
	ProductName string `yaml:"product_name" json:"product_name,omitempty"`
	// Footer line describing who sent the email.
	FooterSentBy string `yaml:"footer_sent_by" json:"footer_sent_by,omitempty"`
	// Copyright line in the email footer.
	Copyright string `yaml:"copyright" json:"copyright,omitempty"`
}

// MailTemplateConfig holds per-locale branding and message overrides from mail_template.yaml.
// Implements mail.MailTemplateStore for use by the renderer.
type MailTemplateConfig struct {
	Locales map[string]map[string]interface{} `yaml:"-"` // set by LoadMailTemplateFile; key = lang
}

type Service struct {
	// `name` determines the name of the service.
	// This value is used, e.g. in the subject header of outgoing emails.
	Name string `yaml:"name" json:"name,omitempty" koanf:"name"`
	// `default_mail_locale` is used when Accept-Language is empty (e.g. API-only calls).
	// Should match ezauth-admin default locale (e.g. "en", "ko").
	DefaultMailLocale string `yaml:"default_mail_locale" json:"default_mail_locale,omitempty" koanf:"default_mail_locale" split_words:"true"`
	// `mail_template_file` path to mail_template.yaml (제품명/푸터/저작권 다국어). 비우면 로케일 파일 사용.
	MailTemplateFile string `yaml:"mail_template_file" json:"mail_template_file,omitempty" koanf:"mail_template_file" split_words:"true"`
}

func (s *Service) Validate() error {
	if len(strings.TrimSpace(s.Name)) == 0 {
		return errors.New("field name must not be empty")
	}
	return nil
}
