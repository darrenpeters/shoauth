# shoauth

`shoauth` is an HTTP package for Go that implements the Shopify oauth2 API.


### Features

* Performs Shopify app installation (displaying auth page and retrieving access tokens)
* Optionally performs webhook installation upon app installation
* Performs Shopify request verification (both user & webhook)
* Acts as typical Go middleware - meaning it can be placed in any middleware stack as long as it is made up of `http.Handler`
* Allows custom failure handlers. 
* Has no dependencies external to the Go standard library.

### Example
```go
package main

import (
	"fmt"
	"github.com/darrenpeters/shoauth"
	"net/http"
	"os"
)

type dummyPersistence struct {}

func (p *dummyPersistence) InstallationExists(shopID string) bool {
	return false // Use any ORM or SQL package directly
}

func (p *dummyPersistence) CreateInstallation(shopID string, accessToken string) error {
	return nil // Use any ORM or SQL package directly
}

func (p *dummyPersistence) HasValidSession(r *http.Request) bool {
	return false // Use gorilla sessions or whatever here
}

func oauthConfig(s *shoauth.ShopifyConfig) {
	s.ClientID = os.Getenv("SHOPIFY_CLIENT_ID")
	s.SharedSecret = os.Getenv("SHOPIFY_SHARED_SECRET")
	s.IsEmbedded = true
	s.Scopes = []string{"read_products", "write_products", "read_customers", "write_customers", "read_orders", "write_orders", "read_shipping", "write_shipping"}
	s.Webhooks = make(map[string]string)
	s.Webhooks["orders/fulfilled"] = "https://yourapp.com/orders/fulfilled"
	s.Webhooks["app/uninstalled"] = "https://yourapp.com/app/uninstalled"
}

func main() {
	shoauthHandler := shoauth.NewShopifyOauthHandler(http.DefaultServeMux, shoauth.DefaultFailureHandler(), &dummyPersistence{}, oauthConfig)
	fmt.Println("Listening on http://127.0.0.1:8000/")
	http.ListenAndServe(":8000", shoauthHandler)
}
```

### Documentation

Documentation is available at [GoDoc](http://godoc.org/github.com/darrenpeters/shoauth)