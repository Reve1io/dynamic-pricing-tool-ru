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

type GetchipsClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewGetchipsClient(baseURL, token string) *GetchipsClient {
	return &GetchipsClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *GetchipsClient) SearchPart(ctx context.Context, partNumber string, quantity int) (*types.GetchipsResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	q := u.Query()
	q.Set("input", partNumber)
	q.Set("qty", strconv.Itoa(quantity))
	q.Set("token", c.token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "PartAPIProcessor/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result types.GetchipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *GetchipsClient) SearchPartAsync(ctx context.Context, partNumber string, quantity int, results chan<- types.APIResponse) {
	go func() {
		data, err := c.SearchPart(ctx, partNumber, quantity)

		select {
		case <-ctx.Done():
			return
		case results <- types.APIResponse{
			PartNumber:   partNumber,
			GetchipsData: data,
			GetchipsErr:  err,
		}:
		}
	}()
}
