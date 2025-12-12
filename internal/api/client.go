package api

import (
	"context"
	"dynamic-pricing-tool-ru/internal/types"
)

type APIClient interface {
	SearchPart(ctx context.Context, partNumber string, quantity int) (*types.GetchipsResponse, error)
	SearchPartAsync(ctx context.Context, partNumber string, quantity int, results chan<- types.APIResponse)
}
