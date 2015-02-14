package shoauth

import "net/http"

type ShopifyPersistence interface {
	InstallationExists(shopID string) bool
	CreateInstallation(shopID string, accessToken string) error
	HasValidSession(r *http.Request) bool
}
