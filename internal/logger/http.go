package logger

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LoggingRoundTripper struct {
	rt http.RoundTripper
}

func NewLoggingRoundTripper(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &LoggingRoundTripper{rt: rt}
}

func sanitizeHeaders(h http.Header) http.Header {
	safe := h.Clone()

	// стандартные варианты авторизации
	safe.Del("Authorization")
	safe.Del("Proxy-Authorization")

	// частые кастомные ключи
	safe.Del("X-API-Key")
	safe.Del("X-Auth-Token")
	safe.Del("Api-Key")
	safe.Del("Token")

	return safe
}

func (l *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// --- Request body ---
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	resp, err := l.rt.RoundTrip(req)
	latency := time.Since(start)

	if err != nil {
		L.Error("API request failed",
			zap.String("method", req.Method),
			zap.String("url", req.URL.String()),
			zap.Any("headers", sanitizeHeaders(req.Header)),
			zap.Duration("latency", latency),
			zap.Error(err),
		)
		return nil, err
	}

	// --- Response body ---
	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)

		// 🔥 ВОТ ЭТО ТЫ ХОТЕЛ
		L.Info("RAW API RESPONSE",
			zap.String("url", req.URL.String()),
			zap.ByteString("raw_body", limitSize(respBody)),
		)

		// возвращаем body обратно, иначе дальше его никто не сможет читать
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	L.Info("API exchange",
		zap.String("method", req.Method),
		zap.String("url", req.URL.String()),
		zap.Any("headers", sanitizeHeaders(req.Header)),
		zap.Int("status", resp.StatusCode),
		zap.Duration("latency", latency),
		zap.ByteString("request_body", limitSize(reqBody)),
		zap.ByteString("response_body", limitSize(respBody)),
	)

	return resp, nil
}

func limitSize(data []byte) []byte {
	const max = 4096
	if len(data) > max {
		return append(data[:max], []byte("...truncated")...)
	}
	return data
}

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
