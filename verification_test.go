package shoauth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const sharedSecret = "abcde"

type requestFromShopify struct {
	Code      string `json:"code"`
	Shop      string `json:"shop"`
	Signature string `json:"signature"`
	Timestamp string `json:"timestamp"`
}

func TestVerificationHMACInHeaderFailure(t *testing.T) {
	requestDataStruct := requestFromShopify{
		Shop:      "some-shop.myshopify.com",
		Code:      "a94a110d86d2452eb3e2af4cfb8a3828",
		Timestamp: "1337178173",
		Signature: "6e39a2ea9e497af6cb806720da1f1bf3",
	}
	requestDataString, _ := json.Marshal(requestDataStruct)
	request, err := http.NewRequest("GET", "/", bytes.NewReader(requestDataString))
	if err != nil {
		t.Errorf("Error creating HTTP request.")
	}
	request.Header.Add("X-Shopify-Hmac-SHA256", "2cb1a277650a659f1b11e92a4a64275b128e037f2c3390e3c8fd2d8721dac9e2")
	if err := validateRequest(request, sharedSecret); err != ErrInvalidHMAC {
		t.Errorf("Request should have failed but was valid.")
	}
}

func TestVerificationHMACInURLFailure(t *testing.T) {
	request, err := http.NewRequest("GET", "/?shop=some-shop.myshopify.com&code=a94a110d86d2452eb3e2af4cfb8a3828&timestamp=1337178173&signature=6e39a2ea9e497af6cb806720da1f1bf3&hmac=2cb1a277650a659f1b11e92a4a64275b128e037f2c3390e3c8fd2d8721dac9e2", nil)
	if err != nil {
		t.Errorf("Error creating HTTP request.")
	}
	if err := validateRequest(request, sharedSecret); err != ErrInvalidHMAC {
		t.Errorf("Request should have failed but was valid.")
	}
}

func TestVerificationHMACInHeaderSuccess(t *testing.T) {
	requestDataStruct := requestFromShopify{
		Shop:      "some-shop.myshopify.com",
		Code:      "a94a110d86d2452eb3e2af4cfb8a3828",
		Timestamp: "1337178173",
		Signature: "6e39a2ea9e497af6cb806720da1f1bf3",
	}
	requestDataString, _ := json.Marshal(requestDataStruct)
	request, err := http.NewRequest("GET", "/", bytes.NewReader(requestDataString))
	if err != nil {
		t.Errorf("Error creating HTTP request.")
	}
	request.Header.Add("X-Shopify-Hmac-SHA256", "RGWZUyrIUA8dN5gN0TsLysYijP4yCHDICx52AZ+0K40=")
	request.Header.Set("Content-type", "application/json")
	if err := validateRequest(request, sharedSecret); err != nil {
		t.Errorf("Request should have succeeded but was invalid.")
	}
}

func TestVerificationHMACInURLSuccess(t *testing.T) {
	request, err := http.NewRequest("GET", "/?shop=some-shop.myshopify.com&code=a94a110d86d2452eb3e2af4cfb8a3828&timestamp=1337178173&signature=6e39a2ea9e497af6cb806720da1f1bf3&hmac=af80820c6de5d631beb8f19b87f4bb35c37b524e79dee76edbb93b2c90cd3d1f", nil)
	if err != nil {
		t.Errorf("Error creating HTTP request.")
	}
	if err := validateRequest(request, sharedSecret); err != nil {
		t.Errorf("Request should have succeeded but was invalid.")
	}
}
