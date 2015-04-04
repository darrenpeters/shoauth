package shoauth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type storeScriptTagRequestContainer struct {
	ScriptTag storeScriptTagRequest `json:"script_tag"`
}

type storeScriptTagRequest struct {
	Event  string `json:"event"`
	Source string `json:"src"`
}

func createScriptTag(shop, accessToken, event, source string) error {
	client := http.Client{}
	var requestData storeScriptTagRequestContainer
	requestData.ScriptTag = storeScriptTagRequest{
		Event:  event,
		Source: source,
	}
	requestDataString, err := json.Marshal(requestData)
	if err != nil {
		return ErrInvalidRequestData
	}
	req, err := http.NewRequest("POST", "https://"+shop+"/admin/script_tags.json", bytes.NewReader(requestDataString))
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
