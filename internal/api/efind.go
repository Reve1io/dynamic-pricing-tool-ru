package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"dynamic-pricing-tool-ru/internal/types"
)

type EfindClient struct {
	baseURL     string
	accessToken string
	client      *http.Client
}

func NewEfindClient(baseURL, accessToken string) *EfindClient {
	return &EfindClient{
		baseURL:     baseURL,
		accessToken: accessToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *EfindClient) SearchPart(ctx context.Context, partNumber string, quantity int) (*types.EfindResponse, error) {
	searchURL := c.baseURL + "/" + url.PathEscape(partNumber)

	reqURL, err := url.Parse(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search URL: %w", err)
	}

	q := reqURL.Query()
	q.Set("access_token", c.accessToken)
	q.Set("stock", "0")
	q.Set("hp", "1")
	q.Set("cur", "usd")
	q.Set("tm", "2")
	q.Set("qty", strconv.Itoa(quantity))

	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "PartAPIProcessor/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result types.EfindResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
