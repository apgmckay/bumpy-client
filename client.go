package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	_ = iota
	v1
)

type Client struct {
	URL        string
	httpClient http.Client
}

var responses map[string]interface{}

func New(endpoint, timeDurationString string) (Client, error) {
	parsedEndpoint, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return Client{}, err
	}

	timeout, err := time.ParseDuration(timeDurationString)
	if err != nil {
		return Client{}, err
	}

	return Client{
		URL: parsedEndpoint.String(),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c Client) GetMajor(params map[string]string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/v%d/major/%s", c.URL, v1, params["version"])
	endpoint = c.genURLQueryParams(endpoint, params)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error")
	}

	json.Unmarshal(body, &responses)

	return responses["version"].(string), nil
}

func (c Client) GetMinor(params map[string]string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/v%d/minor/%s", c.URL, v1, params["version"])
	endpoint = c.genURLQueryParams(endpoint, params)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error")
	}

	json.Unmarshal(body, &responses)

	return responses["version"].(string), nil
}

func (c Client) GetPatch(params map[string]string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/v%d/patch/%s", c.URL, v1, params["version"])

	endpoint = c.genURLQueryParams(endpoint, params)

	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error")
	}

	json.Unmarshal(body, &responses)

	return responses["version"].(string), nil
}

func (c Client) genURLQueryParams(endpoint string, queryParams map[string]string) string {
	firstParam := true

	delete(queryParams, "version")

	for k, v := range queryParams {
		if v != "" {
			if firstParam {
				endpoint = fmt.Sprintf("%s?%s=%s", endpoint, k, v)
				firstParam = false
			} else {
				endpoint = fmt.Sprintf("%s&%s=%s", endpoint, k, v)
			}
		}
	}

	return endpoint
}
