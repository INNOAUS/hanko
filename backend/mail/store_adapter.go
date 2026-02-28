package mail

import (
	"github.com/teamhanko/hanko/backend/v2/config"
)

type configStoreAdapter struct {
	c *config.MailTemplateConfig
}

func (a *configStoreAdapter) GetBranding(lang string) (productName, footerSentBy, copyright string) {
	return a.c.GetBranding(lang)
}

func (a *configStoreAdapter) GetMessage(lang, messageID string, data map[string]interface{}) (string, bool) {
	return a.c.GetMessage(lang, messageID, data)
}

// NewMailTemplateStoreFromConfig wraps *config.MailTemplateConfig as MailTemplateStore for use with NewRenderer.
// Returns nil if c is nil.
func NewMailTemplateStoreFromConfig(c *config.MailTemplateConfig) MailTemplateStore {
	if c == nil {
		return nil
	}
	return &configStoreAdapter{c: c}
}
