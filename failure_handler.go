package shoauth

import "net/http"

type failureHandler struct{}

// DefaultFailureHandler returns an HTTP handler that simply displays a 403.
func DefaultFailureHandler() http.Handler {
	return &failureHandler{}
}

func (f *failureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
