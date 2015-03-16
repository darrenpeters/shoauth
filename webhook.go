package shoauth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type storeWebhookRequestContainer struct {
	Webhook storeWebhookRequest `json:"webhook"`
}

type storeWebhookRequest struct {
	Topic   string `json:"topic"`
	Address string `json:"address"`
	Format  string `json:"format"`
}

func createWebhook(shop, accessToken, webhook, address string) error {
	client := http.Client{}
	var requestData storeWebhookRequestContainer
	requestData.Webhook = storeWebhookRequest{
		Topic:   webhook,
		Address: address,
		Format:  "json",
	}
	requestDataString, err := json.Marshal(requestData)
	if err != nil {
		return ErrInvalidRequestData
	}
	req, err := http.NewRequest("POST", "https://"+shop+"/admin/webhooks.json", bytes.NewReader(requestDataString))
	if err != nil {
		return err
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Shopify-Access-Token", accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return &ErrShopifyHTTPRequestFailed{err: err}
	} else if resp.StatusCode != 201 {
		return &ErrShopifyHTTPRequestFailed{statusCode: resp.StatusCode}
	}

	return nil
}
