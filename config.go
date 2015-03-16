package shoauth

// ShopifyConfig is a structure that contains variables specific to the app
// that the developer is creating.
type ShopifyConfig struct {
	ClientID        string            // Your app's API key
	SharedSecret    string            // Your app's shared secret
	InstallationURI string            // The URI the user is redirected to in order to install your app
	CallbackURI     string            // The base application URI the user is typically redirected to
	HelpURI         string            // The URI the user is redirected to in order to view your app help page
	Webhooks        map[string]string // Webhooks that should be created on installation.
}
