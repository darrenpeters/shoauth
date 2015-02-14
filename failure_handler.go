package shoauth

import "net/http"

type failureHandler struct {
	clientId string
}

// DefaultFailureHandler returns an HTTP handler that simply redirects to the
// app's callback URL.  This will prompt Shopify to redirect the user back into
// the app with a signed request which can be validated to establish a proper
// session.  This works well for quirks such as the app preferences page, which
// is not signed and thus will fail without a valid session.
func DefaultFailureHandler(clientId string) http.Handler {
	return &failureHandler{clientId: clientId}
}

func (f *failureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.FormValue("shop")+"/admin/apps/"+f.clientId, http.StatusMovedPermanently)
}
