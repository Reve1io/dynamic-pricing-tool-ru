package processor

import (
	"dynamic-pricing-tool-ru/internal/types"
)

func AnalyzeResults(results []types.CombinedResult) *types.AnalysisResult {
	var analysis types.AnalysisResult

	for _, result := range results {
		analysis.TotalParts++

		if result.Getchips != nil {
			analysis.GetchipsSuccess++
			if result.Getchips.AvailableQty > 0 {
				analysis.GetchipsInStock++
			}
		}

		if result.Efind != nil {
			analysis.EfindSuccess++
			if result.Efind.InStock {
				analysis.EfindInStock++
			}
		}

		if result.Getchips != nil && result.Efind != nil {
			analysis.BothAPIsSuccess++

			// Сравниваем цены (пример бизнес-логики)
			if result.Getchips.Price > 0 && result.Efind.Price > 0 {
				// Конвертируем в одну валюту для сравнения
				priceComparison := comparePrices(result.Getchips, result.Efind)
				analysis.PriceComparisons = append(analysis.PriceComparisons, priceComparison)
			}
		}
	}

	return &analysis
}

func comparePrices(getchips *types.SimplifiedGetchipsData, efind *types.SimplifiedEfindData) types.PriceComparison {
	comparison := types.PriceComparison{
		PartNumber: getchips.PartNumber,
	}

	// Простая логика сравнения
	if getchips.Currency == efind.Currency {
		comparison.GetchipsPrice = getchips.Price
		comparison.EfindPrice = efind.Price

		if getchips.Price < efind.Price {
			comparison.BetterPrice = "getchips"
			comparison.DifferencePercent = (efind.Price - getchips.Price) / efind.Price * 100
		} else if efind.Price < getchips.Price {
			comparison.BetterPrice = "efind"
			comparison.DifferencePercent = (getchips.Price - efind.Price) / getchips.Price * 100
		} else {
			comparison.BetterPrice = "equal"
		}
	}

	return comparison
}

// Дополнительные типы для анализа
type AnalysisResult struct {
	TotalParts       int
	GetchipsSuccess  int
	GetchipsInStock  int
	EfindSuccess     int
	EfindInStock     int
	BothAPIsSuccess  int
	PriceComparisons []PriceComparison
}

type PriceComparison struct {
	PartNumber        string
	GetchipsPrice     float64
	EfindPrice        float64
	BetterPrice       string
	DifferencePercent float64
}
