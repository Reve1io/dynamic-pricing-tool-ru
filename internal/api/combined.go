package api

import (
	"context"
	"sync"

	"dynamic-pricing-tool-ru/internal/types"
)

type CombinedAPIClient struct {
	getchips *GetchipsClient
	efind    *EfindClient
	promelec *PromelecClient
}

func NewCombinedAPIClient(getchips *GetchipsClient, efind *EfindClient, promelec *PromelecClient) *CombinedAPIClient {
	return &CombinedAPIClient{
		getchips: getchips,
		efind:    efind,
		promelec: promelec,
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

	wg.Add(3)

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

	go func() {
		defer wg.Done()
		data, err := c.promelec.SearchPart(ctx, partNumber)
		result.PromelecData = data
		result.PromelecErr = err
	}()

	wg.Wait()
	return result
}
