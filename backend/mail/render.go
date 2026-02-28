package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed templates/* locales/*
var mailFS embed.FS

// MailTemplateStore provides overrides for mail layout and body from mail_template.yaml.
// GetMessage may return (string, true) to use that string (with template substitution); (_, false) = use bundle.
type MailTemplateStore interface {
	GetBranding(lang string) (productName, footerSentBy, copyright string)
	GetMessage(lang, messageID string, data map[string]interface{}) (string, bool)
}

type Renderer struct {
	templatePlain *template.Template
	bundle        *i18n.Bundle
	localizer     *i18n.Localizer
	store         MailTemplateStore
}

// NewRenderer creates a Renderer. store may be nil (then only bundle/locales are used).
func NewRenderer(store MailTemplateStore) (*Renderer, error) {
	r := &Renderer{store: store}
	bundle := i18n.NewBundle(language.English)
	dir, err := mailFS.ReadDir("locales")
	if err != nil {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	for _, entry := range dir {
		_, _ = bundle.LoadMessageFileFS(mailFS, fmt.Sprintf("locales/%s", entry.Name()))
	}
	r.bundle = bundle

	// add the translate function to the template, so it can be used inside
	funcMap := template.FuncMap{"t": r.translate}
	t := template.New("root").Funcs(funcMap)
	_, err = t.ParseFS(mailFS, "templates/*.txt.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}
	r.templatePlain = t

	return r, nil
}

// preferredTagsForBundle returns language tags for the bundle; adds zh-CN fallback when lang is zh so Chinese is used.
func preferredTagsForBundle(lang string) []string {
	if lang == "zh" {
		return []string{"zh", "zh-CN"}
	}
	return []string{lang}
}

// translate is a helper function to translate texts in a template.
// If r.store provides an override for (lang, messageID), that is used (with template substitution); else the bundle is used.
func (r *Renderer) translate(messageID string, templateData map[string]interface{}) string {
	lang, _ := templateData["renderer_lang"].(string)
	if r.store != nil {
		if s, ok := r.store.GetMessage(lang, messageID, templateData); ok {
			return s
		}
	}
	tags := preferredTagsForBundle(lang)
	localizer := i18n.NewLocalizer(r.bundle, tags...)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
}

// RenderPlain renders a template with the given data and lang.
// The lang can be the contents of Accept-Language headers as defined in http://www.ietf.org/rfc/rfc2616.txt.
func (r *Renderer) RenderPlain(templateName string, lang string, data map[string]interface{}) (string, error) {
	tags := preferredTagsForBundle(lang)
	r.localizer = i18n.NewLocalizer(r.bundle, tags...)
	data["renderer_lang"] = lang
	templateBuffer := &bytes.Buffer{}
	err := r.templatePlain.ExecuteTemplate(templateBuffer, fmt.Sprintf("%s.txt.tmpl", templateName), data)
	if err != nil {
		return "", fmt.Errorf("failed to fill plain text template with data: %w", err)
	}
	return strings.TrimSpace(templateBuffer.String()), nil
}

// RenderHTML renders an HTML template with the given data and lang.
// The lang can be the contents of Accept-Language headers as defined in http://www.ietf.org/rfc/rfc2616.txt.
func (r *Renderer) RenderHTML(templateName string, lang string, data map[string]interface{}) (string, error) {
	var buffer bytes.Buffer

	tags := preferredTagsForBundle(lang)
	r.localizer = i18n.NewLocalizer(r.bundle, tags...)
	data["renderer_lang"] = lang

	templateHTML := template.New("root").Funcs(template.FuncMap{"t": r.translate})
	patterns := []string{"templates/layout.html.tmpl", fmt.Sprintf("templates/%s.html.tmpl", templateName)}
	_, err := templateHTML.ParseFS(mailFS, patterns...)
	if err != nil {
		return "", fmt.Errorf("failed to parse html template: %w", err)
	}

	err = templateHTML.ExecuteTemplate(&buffer, "layout", data)
	if err != nil {
		return "", fmt.Errorf("failed to execute html template: %w", err)
	}

	return strings.TrimSpace(buffer.String()), nil
}

func (r *Renderer) Translate(lang string, messageID string, data map[string]interface{}) string {
	if r.store != nil {
		if s, ok := r.store.GetMessage(lang, messageID, data); ok {
			return s
		}
	}
	tags := preferredTagsForBundle(lang)
	loc := i18n.NewLocalizer(r.bundle, tags...)
	return loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
}
