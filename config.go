package shoauth

// ShopifyConfig is a structure that contains variables specific to the app
// that the developer is creating.
type ShopifyConfig struct {
	ClientID     string              // Your app's API key
	SharedSecret string              // Your app's shared secret
	RedirectURI  string              // If you want a different URL other than the callback on installation, put it here
	HelpURI      string              // The URI the user is redirected to in order to view your app help page
	Scopes       []string            // Your app's required scopes
	IsEmbedded   bool                // If your app is embedded
	Webhooks     map[string]string   // Webhooks that should be created on installation.
	Scripts      map[string][]string // Script tags that should be created on installation.  Map is map[event][]sources
}
