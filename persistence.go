package shoauth

import "net/http"

// ShopifyPersistence contains functions that a client app must implement to
// help this middleware.
type ShopifyPersistence interface {
	InstallationExists(shopID string) bool                      // Checks if a shop has installed an app
	CreateInstallation(shopID string, accessToken string) error // Stores the installation in your app's persistence
	HasValidSession(r *http.Request) bool                       // Tells the middleware if a valid session exists
}
