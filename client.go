package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func (c Client) PostBumpMajor(params map[string]string, body io.Reader) (string, error) {
	return c.do("POST", fmt.Sprintf("bump/major/%s", params["version"]), params, body)
}

func (c Client) PostBumpMinor(params map[string]string, body io.Reader) (string, error) {
	return c.do("POST", fmt.Sprintf("bump/minor/%s", params["version"]), params, body)
}

func (c Client) PostBumpPatch(params map[string]string, body io.Reader) (string, error) {
	return c.do("POST", fmt.Sprintf("bump/patch/%s", params["version"]), params, body)
}

func (c Client) GetBumpMajor(params map[string]string) (string, error) {
	return c.do("GET", fmt.Sprintf("bump/major/%s", params["version"]), params, nil)
}

func (c Client) GetBumpMinor(params map[string]string) (string, error) {
	return c.do("GET", fmt.Sprintf("bumpy/minor/%s", params["version"]), params, nil)
}

func (c Client) GetBumpPatch(params map[string]string) (string, error) {
	return c.do("GET", fmt.Sprintf("bump/patch/%s", params["version"]), params, nil)
}

func (c Client) GetBlocked() (bool, error) {
	result, err := c.do("GET", "blocked", map[string]string{}, nil)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(result)
	if err != nil {
		return false, err
	}
	return b, err
}

func (c Client) do(method, segment string, params map[string]string, body io.Reader) (string, error) {
	endpoint := fmt.Sprintf("%s/api/v%d/%s", c.URL, v1, segment)
	endpoint = c.genURLQueryParams(endpoint, params)

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return "", err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result["version"].(string), nil
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
