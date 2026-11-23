package emails

import "strings"

// URL returns an absolute link to a page of this installation's web app
// (APP_URL), so emails never point to the hosted httpsms.com service.
func (config *HermesGeneratorConfig) URL(path string) string {
	return strings.TrimSuffix(config.AppURL, "/") + path
}
