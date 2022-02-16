package symbiosis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.symbiosis.host/"
	userAgent      = "symbiosis-client-lib/symbiosis-go-api"
	mediaType      = "application/json"
)

type Client struct {
	client *http.Client

	apiKey    string
	BaseURL   *url.URL
	UserAgent string

	Clusters ClusterService
}

type SymbiosisApiError struct {
	StatusCode int
	ErrorType  string `json:"error"`
	Message    string `json:"message"`
	Path       string `json:"path"`
}

func (error *SymbiosisApiError) Error() string {
	return fmt.Sprintf("Symbiosis API Error: %v (type %v, status %v)", error.Message, error.ErrorType, error.StatusCode)
}

func NewClient(httpClient *http.Client, baseURL string, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	parsedBaseURL, _ := url.Parse(baseURL)

	c := &Client{client: httpClient, BaseURL: parsedBaseURL, UserAgent: userAgent, apiKey: apiKey}

	c.Clusters = &ClusterServiceOp{client: c}

	return c
}

func (c *Client) NewRequest(ctx context.Context, method string, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("X-Auth-ApiKey", c.apiKey)
	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		data, err := ioutil.ReadAll(resp.Body)
		errorResponse := &SymbiosisApiError{
			StatusCode: resp.StatusCode,
		}
		if err == nil && len(data) > 0 {
			err := json.Unmarshal(data, errorResponse)
			if err != nil {
				return nil, errors.New(string(data))
			}
		}

		return nil, errorResponse
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}
