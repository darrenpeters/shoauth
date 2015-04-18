package shoauth

// ShopifyPersistence contains functions that a client app must implement to
// help this middleware.
type ShopifyPersistence interface {
	InstallationExists(shopID string) bool                      // Checks if a shop has installed an app
	CreateInstallation(shopID string, accessToken string) error // Stores the installation in your app's persistence
}
