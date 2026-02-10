package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dynamic-pricing-tool-ru/internal/logger"
	"dynamic-pricing-tool-ru/internal/types"

	"go.uber.org/zap"
)

type PromelecClient struct {
	httpClient *http.Client
	url        string
	login      string
	password   string
}

func NewPromelecClient(url, login, password string) *PromelecClient {
	return &PromelecClient{
		httpClient: &http.Client{Timeout: 15 * time.Second},
		url:        "https://aaa.na4u.ru/rpc/",
		login:      login,
		password:   password,
	}
}

type promelecRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Method   string `json:"method"`
	Name     string `json:"name"`
}

func (c *PromelecClient) SearchPart(ctx context.Context, partNumber string) (types.PromelecResponse, error) {
	reqBody := promelecRequest{
		Login:    c.login,
		Password: c.password,
		Method:   "items_data_find",
		Name:     partNumber,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.url,
		bytes.NewReader(bodyBytes),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(len(bodyBytes))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	logger.L.Info("PROMELEC RAW",
		zap.ByteString("raw", raw),
	)

	var result types.PromelecResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		logger.L.Error("PROMELEC DECODE ERROR",
			zap.ByteString("raw", raw),
			zap.Error(err),
		)
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return result, nil
}
