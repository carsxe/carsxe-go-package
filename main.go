package carsxe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a minimal CarsXE API client that works with simple key/value maps.
type Client struct {
	apiKey     string
	baseURL    string
	source     string
	httpClient *http.Client
}

// Option configures a Client instance.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") }
}

// WithHTTPClient allows providing a custom *http.Client.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) { c.httpClient = h }
}

// WithSource changes the default "source" query parameter (default: "go").
func WithSource(src string) Option {
	return func(c *Client) { c.source = src }
}

// New creates a new CarsXE client.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: "https://api.carsxe.com",
		source:  "go",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// buildURL builds a full URL with provided raw map params (no reflection).
func (c *Client) buildURL(endpoint string, params map[string]string) string {
	u, err := url.Parse(c.baseURL + "/" + strings.TrimLeft(endpoint, "/"))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse URL: %v", err))
	}
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("source", c.source)
	for k, v := range params {
		if v != "" {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// doRequest executes the HTTP request and decodes JSON into a generic map.
func (c *Client) doRequest(req *http.Request) map[string]any {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("HTTP request failed: %v", err))
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response body: %v", err))
	}

	if len(bodyBytes) == 0 {
		return map[string]any{}
	}

	var out map[string]any
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		panic(fmt.Sprintf("Failed to decode JSON: %v (body=%s)", err, string(bodyBytes)))
	}
	return out
}

// Get performs a generic GET request to any endpoint with query params.
func (c *Client) Get(endpoint string, params map[string]string) map[string]any {
	if params == nil {
		params = map[string]string{}
	}
	urlStr := c.buildURL(endpoint, params)
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}
	return c.doRequest(req)
}

// postJSON performs a POST with a JSON body (used for image-based endpoints).
func (c *Client) postJSON(endpoint string, body any) map[string]any {
	urlStr := c.buildURL(endpoint, nil)
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			panic(fmt.Sprintf("Failed to encode JSON body: %v", err))
		}
	}
	req, err := http.NewRequest(http.MethodPost, urlStr, &buf)
	if err != nil {
		panic(fmt.Sprintf("Failed to create request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

/*
Convenience methods mirroring the TypeScript SDK.
All of these simply pass through to Get() or postJSON() with map[string]string.
You can remove these if you prefer only the generic Get().
*/

// Specs => GET /specs (vin required; deepdata, disableIntVINDecoding optional)
func (c *Client) Specs(params map[string]string) map[string]any {
	return c.Get("specs", params)
}

// MarketValue => GET /v2/marketvalue (vin)
func (c *Client) MarketValue(params map[string]string) map[string]any {
	return c.Get("v2/marketvalue", params)
}

// History => GET /history (vin)
func (c *Client) History(params map[string]string) map[string]any {
	return c.Get("history", params)
}

// Recalls => GET /v1/recalls (vin)
func (c *Client) Recalls(params map[string]string) map[string]any {
	return c.Get("v1/recalls", params)
}

// InternationalVINDecoder => GET /v1/international-vin-decoder (vin)
func (c *Client) InternationalVINDecoder(params map[string]string) map[string]any {
	return c.Get("v1/international-vin-decoder", params)
}

// PlateDecoder => GET /v2/platedecoder (plate, country, state?, district?)
func (c *Client) PlateDecoder(params map[string]string) map[string]any {
	return c.Get("v2/platedecoder", params)
}

// PlateImageRecognition => POST /platerecognition with JSON {"image": "<url>"}
func (c *Client) PlateImageRecognition(imageURL string) map[string]any {
	if strings.TrimSpace(imageURL) == "" {
		panic("image URL required")
	}
	return c.postJSON("platerecognition", map[string]string{"image": imageURL})
}

// VinOCR => POST /v1/vinocr with JSON {"image": "<url>"}
func (c *Client) VinOCR(imageURL string) map[string]any {
	if strings.TrimSpace(imageURL) == "" {
		panic("image URL required")
	}
	return c.postJSON("v1/vinocr", map[string]string{"image": imageURL})
}

// YearMakeModel => GET /v1/ymm (year, make, model, trim?)
func (c *Client) YearMakeModel(params map[string]string) map[string]any {
	return c.Get("v1/ymm", params)
}

// Images => GET /images (make, model, optional year, trim, color, etc.)
func (c *Client) Images(params map[string]string) map[string]any {
	return c.Get("images", params)
}

// ObdCodesDecoder => GET /obdcodesdecoder (code)
func (c *Client) ObdCodesDecoder(params map[string]string) map[string]any {
	return c.Get("obdcodesdecoder", params)
}