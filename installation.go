package shoauth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type storeOauthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type storeOauthResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *shopifyOauthHandler) performInstallation(shop, code string) error {
	if s.InstallationExists(shop) {
		return ErrInstallationExists
	}

	client := http.Client{}
	var requestData storeOauthRequest
	requestData.ClientID = s.config.ClientID
	requestData.ClientSecret = s.config.SharedSecret
	requestData.Code = code
	requestDataString, err := json.Marshal(requestData)
	if err != nil {
		return ErrInvalidRequestData
	}
	resp, err := client.Post("https://"+shop+"/admin/oauth/access_token", "application/json", bytes.NewReader(requestDataString))
	if err != nil {
		return &ErrShopifyHTTPRequestFailed{err: err}
	} else if resp.StatusCode != 200 {
		return &ErrShopifyHTTPRequestFailed{statusCode: resp.StatusCode}
	}
	defer resp.Body.Close()

	var responseData storeOauthResponse
	if err = json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return ErrInvalidResponseData
	}

	if err = s.ShopifyPersistence.CreateInstallation(shop, responseData.AccessToken); err != nil {
		return ErrBadPersistence
	}

	// Create our webhooks
	for webhook, address := range s.config.Webhooks {
		if err = createWebhook(shop, responseData.AccessToken, webhook, address); err != nil {
			return &ErrShopifyHTTPRequestFailed{err: err}
		}
	}

	return nil
}
