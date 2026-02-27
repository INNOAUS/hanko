package config

import (
	"errors"
	"strings"
)

type Service struct {
	// `name` determines the name of the service.
	// This value is used, e.g. in the subject header of outgoing emails.
	Name string `yaml:"name" json:"name,omitempty" koanf:"name"`
	// `default_mail_locale` is used when Accept-Language is empty (e.g. API-only calls).
	// Should match ezauth-admin default locale (e.g. "en", "ko").
	DefaultMailLocale string `yaml:"default_mail_locale" json:"default_mail_locale,omitempty" koanf:"default_mail_locale" split_words:"true"`
}

func (s *Service) Validate() error {
	if len(strings.TrimSpace(s.Name)) == 0 {
		return errors.New("field name must not be empty")
	}
	return nil
}
