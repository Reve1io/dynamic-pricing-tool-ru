package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"dynamic-pricing-tool-ru/internal/logger"
	"dynamic-pricing-tool-ru/internal/types"

	"github.com/yosuke-furukawa/json5/encoding/json5"
	"go.uber.org/zap"
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
			Transport: logger.NewLoggingRoundTripper(nil),
			Timeout:   30 * time.Second,
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

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	trim := bytes.TrimSpace(body)

	if len(trim) == 0 || trim[0] != '[' {
		return nil, fmt.Errorf("efind returned non-json response: %s", string(trim[:200]))
	}

	var result types.EfindResponse
	if err := json5.Unmarshal(trim, &result); err != nil {
		logger.L.Error("DECODE ERROR",
			zap.ByteString("raw", trim),

			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
