package api

import (
	"context"
	"sync"

	"dynamic-pricing-tool-ru/internal/types"
)

type CombinedAPIClient struct {
	getchips *GetchipsClient
	efind    *EfindClient
}

func NewCombinedAPIClient(getchips *GetchipsClient, efind *EfindClient) *CombinedAPIClient {
	return &CombinedAPIClient{
		getchips: getchips,
		efind:    efind,
	}
}

func (c *CombinedAPIClient) SearchPart(ctx context.Context, partNumber string, quantity int) (*types.APIResponse, error) {
	return nil, nil
}

func (c *CombinedAPIClient) SearchAllAPIs(ctx context.Context, partNumber string, quantity int) types.APIResponse {
	var wg sync.WaitGroup
	var result types.APIResponse

	result.PartNumber = partNumber
	result.RequestedQty = quantity

	wg.Add(2)

	go func() {
		defer wg.Done()
		data, err := c.getchips.SearchPart(ctx, partNumber, quantity)
		result.GetchipsData = data
		result.GetchipsErr = err
	}()

	go func() {
		defer wg.Done()
		data, err := c.efind.SearchPart(ctx, partNumber, quantity)
		result.EfindData = data
		result.EfindErr = err
	}()

	wg.Wait()
	return result
}
