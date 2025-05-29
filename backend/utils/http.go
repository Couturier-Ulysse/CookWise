package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DoOptions envoie une requête HTTP OPTIONS (nécessaire pour certaines APIs comme Jow)
func DoOptions(endpoint string, headers map[string]string, params map[string]string) error {
	fullURL := endpoint + "?" + encodeParams(params)

	req, err := http.NewRequest("OPTIONS", fullURL, nil)
	if err != nil {
		return fmt.Errorf("creating OPTIONS request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("executing OPTIONS request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("OPTIONS request returned status %d", resp.StatusCode)
	}

	return nil
}

// DoPost effectue une requête POST avec headers et paramètres
func DoPost(endpoint string, headers map[string]string, params map[string]string, body []byte) (*http.Response, error) {
	fullURL := endpoint + "?" + encodeParams(params)

	req, err := http.NewRequest("POST", fullURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("creating POST request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing POST request: %w", err)
	}

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("POST request failed: status %d, body: %s", resp.StatusCode, respBody)
	}

	return resp, nil
}

// encodeParams transforme une map en string URL-encodée
func encodeParams(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	return q.Encode()
}
