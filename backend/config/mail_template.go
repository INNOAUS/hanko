package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// LoadMailTemplateFile reads mail_template.yaml and returns a config that provides branding and message overrides per locale.
// Path is relative to process cwd or absolute. If not found, tries several relative paths so that e.g. CWD backend/hanko/backend
// still finds backend/config/mail_template.yaml (../../config/mail_template.yaml).
// Returns nil and nil error if file is empty or missing.
func LoadMailTemplateFile(path string) (*MailTemplateConfig, error) {
	if path == "" {
		return nil, nil
	}
	tryPaths := []string{path}
	if !filepath.IsAbs(path) {
		base := filepath.Base(path)
		tryPaths = append(tryPaths, filepath.Join("config", base))
		tryPaths = append(tryPaths, filepath.Join("..", "config", base))
		tryPaths = append(tryPaths, filepath.Join("../..", "config", base))
	}
	var data []byte
	var err error
	for _, p := range tryPaths {
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("mail_template_file: %w", err)
		}
	}
	if err != nil || data == nil {
		return nil, nil
	}
	var raw map[string]map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("mail_template_file parse: %w", err)
	}
	if len(raw) == 0 {
		return nil, nil
	}
	return &MailTemplateConfig{Locales: raw}, nil
}

func (m *MailTemplateConfig) getEntry(lang string) map[string]interface{} {
	if m == nil || m.Locales == nil {
		return nil
	}
	if entry, ok := m.Locales[lang]; ok {
		return entry
	}
	for i := 0; i < len(lang); i++ {
		if lang[i] == '-' {
			if entry, ok := m.Locales[lang[:i]]; ok {
				return entry
			}
			break
		}
	}
	return nil
}

func getStr(entry map[string]interface{}, key string) string {
	if entry == nil {
		return ""
	}
	v, ok := entry[key]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

// BrandingForLang returns branding for the given language.
func (m *MailTemplateConfig) BrandingForLang(lang string) MailBranding {
	entry := m.getEntry(lang)
	if entry == nil {
		return MailBranding{}
	}
	return MailBranding{
		ProductName:  getStr(entry, "product_name"),
		FooterSentBy: getStr(entry, "footer_sent_by"),
		Copyright:    getStr(entry, "copyright"),
	}
}

// GetBranding returns product name, footer text, and copyright for the given language.
func (m *MailTemplateConfig) GetBranding(lang string) (productName, footerSentBy, copyright string) {
	b := m.BrandingForLang(lang)
	return b.ProductName, b.FooterSentBy, b.Copyright
}

// GetMessage returns the message string for (lang, messageID) with template substitution applied. ok is false if not found.
func (m *MailTemplateConfig) GetMessage(lang, messageID string, data map[string]interface{}) (string, bool) {
	entry := m.getEntry(lang)
	if entry == nil {
		return "", false
	}
	raw, ok := entry[messageID]
	if !ok || raw == nil {
		return "", false
	}
	s, ok := raw.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return "", false
	}
	// Substitute {{ .Key }} with data
	tpl, err := template.New("msg").Parse(s)
	if err != nil {
		return s, true
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return s, true
	}
	return buf.String(), true
}
