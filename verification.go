package shoauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"sort"
)

func validateRequest(r *http.Request, sharedSecret string) error {

	var hmacString, expectedHMAC string
	hmacString = r.Header.Get("X-Shopify-Hmac-SHA256")
	// If the HMAC string came in the header, we generate the hash from the
	// request body.
	if len(hmacString) > 0 {
		bodyContents, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyContents))
		mac := hmac.New(sha256.New, []byte(sharedSecret))
		mac.Write(bodyContents)
		expectedHMAC = base64.StdEncoding.EncodeToString(mac.Sum(nil))
		// Otherwise, we generate the hash from the request form parameters.
	} else {
		hmacString = r.FormValue("hmac")
		r.Form.Del("hmac")
		r.Form.Del("signature")
		var formKeys []string
		for key := range r.Form {
			formKeys = append(formKeys, key)
		}
		sort.Strings(formKeys)
		hashString := ""
		for i, key := range formKeys {
			if i == 0 {
				hashString += key + "=" + r.FormValue(key)
			} else {
				hashString += "&" + key + "=" + r.FormValue(key)
			}
		}
		mac := hmac.New(sha256.New, []byte(sharedSecret))
		mac.Write([]byte(hashString))
		expectedHMAC = hex.EncodeToString(mac.Sum(nil))
	}

	if !hmac.Equal([]byte(hmacString), []byte(expectedHMAC)) {
		return ErrInvalidHMAC
	}
	return nil
}
