package processor

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"dynamic-pricing-tool-ru/internal/api"
	"dynamic-pricing-tool-ru/internal/types"
)

type Processor struct {
	combinedClient *api.CombinedAPIClient
	chunkSize      int
	workerPoolSize int
}

func NewProcessorWithClients(getchipsClient *api.GetchipsClient, efindClient *api.EfindClient, promelec *api.PromelecClient, chunkSize int) *Processor {
	combinedClient := api.NewCombinedAPIClient(getchipsClient, efindClient, promelec)
	return &Processor{
		combinedClient: combinedClient,
		chunkSize:      chunkSize,
		workerPoolSize: 20,
	}
}

func (p *Processor) ProcessRequest(ctx context.Context, req *types.Request) ([]types.CombinedResult, error) {
	parts, err := p.extractPartData(req)
	if err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return []types.CombinedResult{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan types.PartData, len(parts))
	resultsChan := make(chan types.CombinedResult, len(parts))

	var wg sync.WaitGroup

	for i := 0; i < p.workerPoolSize; i++ {
		wg.Add(1)
		go p.worker(ctx, jobs, resultsChan, &wg)
	}

	go func() {
		for _, part := range parts {
			select {
			case <-ctx.Done():
				return
			case jobs <- part:
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var allResults []types.CombinedResult
	for result := range resultsChan {
		allResults = append(allResults, result)
	}

	return allResults, nil
}

func (p *Processor) worker(ctx context.Context, jobs <-chan types.PartData, results chan<- types.CombinedResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for part := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			qty := p.parseQuantity(part.Quantity)

			apiResult := p.combinedClient.SearchAllAPIs(ctx, part.PartNumber, qty)

			simplifiedGetchips := FormatGetchipsData(apiResult.GetchipsData, part.PartNumber)
			simplifiedEfind := FormatEfindData(apiResult.EfindData)
			promelecFormatted := FormatPromelecData(apiResult.PromelecData)

			combinedResult := types.CombinedResult{
				PartNumber:   part.PartNumber,
				RequestedQty: qty,
				RowIndex:     part.RowIndex,
				Getchips:     simplifiedGetchips,
				Efind:        simplifiedEfind,
				Promelec:     promelecFormatted,
				GetchipsRaw:  apiResult.GetchipsData,
				EfindRaw:     apiResult.EfindData,
				Timestamp:    time.Now().Format(time.RFC3339),
			}

			if apiResult.GetchipsErr != nil {
				combinedResult.GetchipsError = apiResult.GetchipsErr.Error()
			}

			if apiResult.EfindErr != nil {
				combinedResult.EfindError = apiResult.EfindErr.Error()
			}

			if apiResult.PromelecErr != nil {
				combinedResult.PromelecError = apiResult.PromelecErr.Error()
			}

			// Отправляем результат
			select {
			case <-ctx.Done():
				return
			case results <- combinedResult:
			}
		}
	}
}

func (p *Processor) parseQuantity(quantityStr string) int {
	if quantityStr == "" {
		return 1
	}

	quantityStr = strings.TrimSpace(quantityStr)

	qty, err := strconv.Atoi(quantityStr)
	if err != nil {
		return 1
	}

	if qty <= 0 {
		return 1
	}

	return qty
}

func (p *Processor) extractPartData(req *types.Request) ([]types.PartData, error) {
	if len(req.Data) < 2 {
		return nil, fmt.Errorf("insufficient data")
	}

	var parts []types.PartData

	partNumberIndex := -1
	quantityIndex := -1

	for key, value := range req.Mapping {
		switch value {
		case "partNumber":
			if idx, err := strconv.Atoi(key); err == nil {
				partNumberIndex = idx
			}
		case "quantity":
			if idx, err := strconv.Atoi(key); err == nil {
				quantityIndex = idx
			}
		}
	}

	if partNumberIndex == -1 {
		return nil, fmt.Errorf("partNumber mapping not found")
	}

	for i := 1; i < len(req.Data); i++ {
		row := req.Data[i]

		if len(row) <= partNumberIndex {
			continue
		}

		partNumber := strings.TrimSpace(row[partNumberIndex])
		if partNumber == "" {
			continue
		}

		quantity := ""
		if quantityIndex != -1 && len(row) > quantityIndex {
			quantity = strings.TrimSpace(row[quantityIndex])
		}

		parts = append(parts, types.PartData{
			PartNumber: partNumber,
			Quantity:   quantity,
			RowIndex:   i,
		})
	}

	return parts, nil
}
