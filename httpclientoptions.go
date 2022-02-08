package kraken

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

// HTTPClientOption options used when creating a new HTTPClient
type HTTPClientOption func(c *HTTPClient) error

// HTTPClientWithHTTPClient set the http client of the Kraken client wrapper
func HTTPClientWithHTTPClient(httpClient *http.Client) HTTPClientOption {
	return HTTPClientOption(func(c *HTTPClient) error {
		c.httpClient = httpClient

		return nil
	})
}

// HTTPClientWithBaseURL set the base url of the Kraken client wrapper
func HTTPClientWithBaseURL(baseURL string) HTTPClientOption {
	return HTTPClientOption(func(c *HTTPClient) error {
		if _, err := url.Parse(c.baseURL); err != nil {
			return err
		}

		c.baseURL = baseURL

		return nil
	})
}

// HTTPClientWithSecret set the secret of the Kraken client wrapper
func HTTPClientWithSecret(secret string) HTTPClientOption {
	return HTTPClientOption(func(c *HTTPClient) error {
		if _, err := base64.StdEncoding.DecodeString(c.secret); err != nil {
			return fmt.Errorf("invalid secret: %s", err)
		}

		c.secret = secret

		return nil
	})
}

// HTTPClientDryRun set the Kraken client to not execute requests
func HTTPClientDryRun() HTTPClientOption {
	return HTTPClientOption(func(c *HTTPClient) error {
		c.dryRun = true

		return nil
	})
}
