package carsxe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
func (c *Client) buildURL(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL + "/" + strings.TrimLeft(endpoint, "/"))
	if err != nil {
		return "", err
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
	return u.String(), nil
}

// doRequest executes the HTTP request and decodes JSON into a generic map.
func (c *Client) doRequest(req *http.Request) (map[string]any, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("carsxe: non-2xx response (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	if len(bodyBytes) == 0 {
		return map[string]any{}, nil
	}

	var out map[string]any
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return nil, fmt.Errorf("carsxe: decode error: %w (body=%s)", err, string(bodyBytes))
	}
	return out, nil
}

// Get performs a generic GET request to any endpoint with query params.
func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string) (map[string]any, error) {
	if params == nil {
		params = map[string]string{}
	}
	urlStr, err := c.buildURL(endpoint, params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// postJSON performs a POST with a JSON body (used for image-based endpoints).
func (c *Client) postJSON(ctx context.Context, endpoint string, body any) (map[string]any, error) {
	urlStr, err := c.buildURL(endpoint, nil)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlStr, &buf)
	if err != nil {
		return nil, err
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
func (c *Client) Specs(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "specs", params)
}

// MarketValue => GET /v2/marketvalue (vin)
func (c *Client) MarketValue(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "v2/marketvalue", params)
}

// History => GET /history (vin)
func (c *Client) History(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "history", params)
}

// Recalls => GET /v1/recalls (vin)
func (c *Client) Recalls(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "v1/recalls", params)
}

// InternationalVINDecoder => GET /v1/international-vin-decoder (vin)
func (c *Client) InternationalVINDecoder(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "v1/international-vin-decoder", params)
}

// PlateDecoder => GET /v2/platedecoder (plate, country, state?, district?)
func (c *Client) PlateDecoder(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "v2/platedecoder", params)
}

// PlateImageRecognition => POST /platerecognition with JSON {"image": "<url>"}
func (c *Client) PlateImageRecognition(ctx context.Context, imageURL string) (map[string]any, error) {
	if strings.TrimSpace(imageURL) == "" {
		return nil, errors.New("image URL required")
	}
	return c.postJSON(ctx, "platerecognition", map[string]string{"image": imageURL})
}

// VinOCR => POST /v1/vinocr with JSON {"image": "<url>"}
func (c *Client) VinOCR(ctx context.Context, imageURL string) (map[string]any, error) {
	if strings.TrimSpace(imageURL) == "" {
		return nil, errors.New("image URL required")
	}
	return c.postJSON(ctx, "v1/vinocr", map[string]string{"image": imageURL})
}

// YearMakeModel => GET /v1/ymm (year, make, model, trim?)
func (c *Client) YearMakeModel(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "v1/ymm", params)
}

// Images => GET /images (make, model, optional year, trim, color, etc.)
func (c *Client) Images(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "images", params)
}

// ObdCodesDecoder => GET /obdcodesdecoder (code)
func (c *Client) ObdCodesDecoder(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.Get(ctx, "obdcodesdecoder", params)
}