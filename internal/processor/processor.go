package processor

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"dynamic-pricing-tool-ru/internal/api"
	"dynamic-pricing-tool-ru/internal/types"
)

type Processor struct {
	apiClient      api.APIClient
	chunkSize      int
	workerPoolSize int
}

func NewProcessor(apiClient api.APIClient, chunkSize int) *Processor {
	return &Processor{
		apiClient:      apiClient,
		chunkSize:      chunkSize,
		workerPoolSize: 20,
	}
}

func (p *Processor) ProcessRequest(ctx context.Context, req *types.Request) ([]types.APIResponse, error) {
	// Extract part data from request
	parts, err := p.extractPartData(req)
	if err != nil {
		return nil, err
	}

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Создаем каналы
	jobs := make(chan types.PartData, len(parts))
	resultsChan := make(chan types.APIResponse, len(parts))

	var wg sync.WaitGroup

	// Запускаем воркеры
	for w := 0; w < p.workerPoolSize; w++ {
		wg.Add(1)
		go p.worker(ctx, jobs, resultsChan, &wg)
	}

	// Отправляем задания в канал
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

	// Закрываем канал результатов после завершения всех воркеров
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Собираем результаты
	var allResults []types.APIResponse
	for result := range resultsChan {
		allResults = append(allResults, result)
	}

	return allResults, nil
}

func (p *Processor) worker(ctx context.Context, jobs <-chan types.PartData, results chan<- types.APIResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	for part := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			qty, _ := strconv.Atoi(part.Quantity)
			if qty == 0 {
				qty = 1
			}

			// Синхронный вызов для упрощения
			data, err := p.apiClient.SearchPart(ctx, part.PartNumber, qty)

			select {
			case <-ctx.Done():
				return
			case results <- types.APIResponse{
				PartNumber: part.PartNumber,
				Data:       data,
				Error:      err,
			}:
			}
		}
	}
}

func (p *Processor) extractPartData(req *types.Request) ([]types.PartData, error) {
	if len(req.Data) < 2 {
		return nil, fmt.Errorf("insufficient data")
	}

	var parts []types.PartData

	// Skip header row (index 0)
	for i := 1; i < len(req.Data); i++ {
		row := req.Data[i]
		if len(row) < 2 {
			continue
		}

		part := types.PartData{
			PartNumber: row[0], // MPN
			Quantity:   row[1], // Qty
		}
		parts = append(parts, part)
	}

	return parts, nil
}

func (p *Processor) splitIntoChunks(parts []types.PartData, chunkSize int) [][]types.PartData {
	var chunks [][]types.PartData

	for i := 0; i < len(parts); i += chunkSize {
		end := i + chunkSize
		if end > len(parts) {
			end = len(parts)
		}
		chunks = append(chunks, parts[i:end])
	}

	return chunks
}

func (p *Processor) processChunk(ctx context.Context, chunk []types.PartData, results chan<- types.APIResponse) {
	var wg sync.WaitGroup

	for _, part := range chunk {
		wg.Add(1)
		go func(currentPart types.PartData) {
			defer wg.Done()

			qty, _ := strconv.Atoi(currentPart.Quantity)
			if qty == 0 {
				qty = 1 // Default quantity
			}

			// Используем apiClient из Processor (p.apiClient)
			p.apiClient.SearchPartAsync(ctx, currentPart.PartNumber, qty, results)
		}(part)
	}

	wg.Wait()
}
