package shoauth

import (
	"errors"
	"fmt"
)

var (
	// ErrInstallationExists is returned if a shop that has already installed this
	// application attempts to install it a second time.
	ErrInstallationExists = errors.New("An installation already exists for this shop and application.")
	// ErrInvalidHMAC is returned when a request from shopify (or pretending to
	// be from shopify) contains an invalid HMAC signature.
	ErrInvalidHMAC = errors.New("The request contained an invalid HMAC - this request will be discarded.")
	// ErrInvalidRequestData is returned when a json.Marshal fails on the
	// shopify access token request struct.
	ErrInvalidRequestData = errors.New("The request data to be sent to shopify for an access token was invalid - please make sure your configuration settings are applied properly.")
	// ErrInvalidResponseData is returned when a json.Unmarshal fails on a
	// response from shopify
	ErrInvalidResponseData = errors.New("The response data returned from shopify was in an unexpected format.")
	// ErrBadPersistence is returned when the implementation of persistence
	// passed to shoauth fails for some reason.
	ErrBadPersistence = errors.New("The persistence passed to shoauth failed to perform an action.")
)

// ErrShopifyHTTPRequestFailed is returned when an HTTP request to a
// shopify shop fails.
type ErrShopifyHTTPRequestFailed struct {
	err        error
	statusCode int
}

func (e *ErrShopifyHTTPRequestFailed) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return fmt.Sprintf("There was an error contacting the shopify shop.  The HTTP status code returned was %d.", e.statusCode)
}
