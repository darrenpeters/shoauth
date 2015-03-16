package shoauth

import "net/http"

type shopifyOauthHandler struct {
	successHandler http.Handler
	failureHandler http.Handler
	config         *ShopifyConfig
	ShopifyPersistence
}

// NewShopifyOauthHandler returns the middleware handler that handles Shopify
// oauth requests and responses.  It will call successHandler.ServeHTTP on a
// successful installation or verification, and will call
// failureHandler.ServeHTTP on an unsuccessful installation or verification.
// The user must pass a shopifyPersistence-satisfying struct and any functions
// they wish to operate on the default config object.
func NewShopifyOauthHandler(successHandler http.Handler, failureHandler http.Handler, persistence ShopifyPersistence, configOptions ...func(*ShopifyConfig)) *shopifyOauthHandler {
	// Set some sensible defaults.
	config := &ShopifyConfig{
		InstallationURI: "/",
		CallbackURI:     "/",
		HelpURI:         "/help",
		Webhooks:        make(map[string]string),
	}

	// Apply the custom config functions passed.
	for _, f := range configOptions {
		f(config)
	}

	return &shopifyOauthHandler{
		successHandler:     successHandler,
		failureHandler:     failureHandler,
		ShopifyPersistence: persistence,
		config:             config,
	}
}

func (s *shopifyOauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// If the user has authenticated via the initial app Callback, the app
	// should have registered a valid session for the user.  As long as that
	// session is active, we do not need to validate requests from said user.
	// The help page is also static and unsigned - we can just display it.
	if s.HasValidSession(r) || r.URL.Path == s.config.HelpURI {
		s.successHandler.ServeHTTP(w, r)
		return
	}

	// Installation requests can be confused by the fact that the installation
	// URI can theoretically be the same as the callback URI.  Thus we need to
	// attempt installation in two scenarios:
	// 1. The request URI is equal to the installation URI, and the
	//    installation URI is different than the callback URI - this is the
	//    basic installation scenario.
	// 2. The request URI is equal to the installation URI, and the
	//    installation URI is equal to the callback URI.  In this instance, in
	//    order to confirm that an installation is required, we validate that
	//    an installation does not exist for the shop making the request.
	if (r.URL.Path == s.config.InstallationURI && s.config.InstallationURI != s.config.CallbackURI) ||
		(r.URL.Path == s.config.InstallationURI && s.config.InstallationURI == s.config.CallbackURI && !s.InstallationExists(r.FormValue("shop"))) {
		// We perform the installation - if it fails, call the app's
		// failure handler.  Otherwise, we redirect to the app's callback URI.
		if err := s.performInstallation(r.FormValue("shop"), r.FormValue("code")); err != nil {
			s.failureHandler.ServeHTTP(w, r)
		} else {
			s.successHandler.ServeHTTP(w, r)
		}
		// If this is not an installation request, we must validate that it has
		// actually come from shopify according to their predefined rules.
	} else {
		if err := validateRequest(r, s.config.SharedSecret); err != nil {
			s.failureHandler.ServeHTTP(w, r)
		} else {
			s.successHandler.ServeHTTP(w, r)
		}
	}
}
