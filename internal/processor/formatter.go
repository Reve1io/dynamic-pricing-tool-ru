package processor

import (
	"fmt"
	"strconv"

	"dynamic-pricing-tool-ru/internal/types"
	"dynamic-pricing-tool-ru/internal/utils"
)

func toInt(v interface{}) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}

func FormatGetchipsData(raw *types.GetchipsResponse, requestedMPN string, requestedQty int) []types.UnifiedOffer {
	if raw == nil {
		return nil
	}

	var offers []types.UnifiedOffer

	for _, d := range raw.Data {
		currency := "USD"

		delivery := formatDeliveryTime(d.Orderdays, "2-3 недели")

		pb := make([]types.PriceBreak, 0, len(d.PriceBreak))
		for _, p := range d.PriceBreak {
			pb = append(pb, types.PriceBreak{
				Quantity:     p.Quantity,
				Price:        p.Price,
				DeliveryTime: delivery,
			})
		}

		priceBreaks := buildPriceBreaks(pb, currency)

		basePrice := 0.0
		if len(pb) > 0 {
			basePrice = pb[0].Price
		}

		offers = append(offers, types.UnifiedOffer{
			MPN:          d.Title,
			RequestedMPN: requestedMPN,
			RequestedQty: requestedQty,

			Manufacturer:   d.Brand,
			SellerName:     "Getchips",
			SellerVerified: true,

			Stock:        d.Quantity,
			Status:       "Найдено",
			Price:        basePrice,
			Currency:     currency,
			DeliveryTime: delivery,

			PriceBreaks: priceBreaks,
			Source:      "getchips",
		})
	}

	return offers
}

func FormatEfindData(raw *types.EfindResponse, requestedMPN string, requestedQty int) []types.UnifiedOffer {
	if raw == nil || len(*raw) == 0 {
		return nil
	}

	var offers []types.UnifiedOffer

	for _, stock := range *raw {
		for _, row := range stock.Rows {

			availableQty := toInt(row.Stock)

			// Собираем pricebreaks
			var pbs []types.PriceBreak
			for _, p := range row.Price {
				if len(p) < 3 {
					continue
				}

				qty := toInt(p[0])
				price, ok := p[2].(float64)
				if !ok {
					continue
				}

				pbs = append(pbs, types.PriceBreak{
					Quantity: qty,
					Price:    price,
				})
			}

			priceBreaks := buildPriceBreaks(pbs, row.Cur)

			basePrice := 0.0
			if len(pbs) > 0 {
				basePrice = pbs[0].Price
			}

			offers = append(offers, types.UnifiedOffer{
				MPN:          row.Part,
				RequestedMPN: requestedMPN,
				RequestedQty: requestedQty,

				Manufacturer: "", // у efind нет бренда
				Description:  "",
				ImageURL:     "",

				SellerName:     stock.StockData.Title,
				SellerHomepage: stock.StockData.Site,
				SellerVerified: true,

				Stock:    availableQty,
				Status:   "Найдено",
				Price:    basePrice,
				Currency: row.Cur,

				PriceBreaks: priceBreaks,
				Source:      "efind",
			})
		}
	}

	return offers
}

func FormatPromelecData(data types.PromelecResponse, requestedMPN string, requestedQty int) []types.UnifiedOffer {
	var offers []types.UnifiedOffer

	for _, item := range data {

		// если vendors нет → создаем 1 оффер с дефолтом
		if len(item.Vendors) == 0 {

			delivery := "1-2 недели"

			var priceBreaks []types.PriceBreak
			for _, pb := range item.Pricebreaks {
				priceBreaks = append(priceBreaks, types.PriceBreak{
					Quantity:     pb.Quant,
					Price:        pb.Price,
					DeliveryTime: delivery,
				})
			}

			basePrice := 0.0
			if len(priceBreaks) > 0 {
				basePrice = priceBreaks[0].Price
			}

			offers = append(offers, types.UnifiedOffer{
				MPN:          item.Name,
				RequestedMPN: requestedMPN,
				RequestedQty: requestedQty,
				Manufacturer: item.ProducerName,
				SellerName:   "Promelec",
				Stock:        item.Quant,
				Status:       "Найдено",
				Price:        basePrice,
				Currency:     "RUB",
				PriceBreaks:  buildPriceBreaks(priceBreaks, "RUB"),
				DeliveryTime: delivery,
				Source:       "promelec",
			})

			continue
		}

		// vendors есть → делаем несколько офферов
		for _, v := range item.Vendors {

			delivery := formatDeliveryTime(v.Delivery, "1-2 недели")

			var priceBreaks []types.PriceBreak
			for _, pb := range v.PriceBreaks {
				priceBreaks = append(priceBreaks, types.PriceBreak{
					Quantity:     pb.Quant,
					Price:        pb.Price,
					DeliveryTime: delivery,
				})
			}

			offers = append(offers, types.UnifiedOffer{
				MPN:          item.Name,
				RequestedMPN: requestedMPN,
				RequestedQty: requestedQty,
				Manufacturer: item.ProducerName,
				SellerName:   "Promelec",
				Stock:        v.Quant,
				Status:       "Найдено",
				Price:        v.PriceBreaks[0].Price,
				Currency:     "RUB",
				PriceBreaks:  buildPriceBreaks(priceBreaks, "RUB"),
				DeliveryTime: "1-2 недели", // верхний уровень всегда default
				Source:       "promelec",
			})
		}
	}

	return offers
}

// CompareAndSelectBest сравнивает результаты от двух API и выбирает лучший
func CompareAndSelectBest(getchips *types.SimplifiedGetchipsData, efind *types.SimplifiedEfindData) (*types.SimplifiedGetchipsData, *types.SimplifiedEfindData, string) {
	if getchips == nil && efind == nil {
		return nil, nil, "no_data"
	}

	if getchips == nil {
		return nil, efind, "efind_only"
	}

	if efind == nil {
		return getchips, nil, "getchips_only"
	}

	// Простая логика сравнения
	// 1. Проверяем наличие
	if getchips.AvailableQty > 0 && efind.AvailableQty == 0 {
		return getchips, nil, "getchips_better_stock"
	}

	if efind.AvailableQty > 0 && getchips.AvailableQty == 0 {
		return nil, efind, "efind_better_stock"
	}

	// По умолчанию возвращаем оба
	return getchips, efind, "both_valid"
}

const (
	deliveryCoef = 1.27
	//markup       = 1.18
)

func buildPriceBreaks(priceBreaks []types.PriceBreak, currency string) []types.UnifiedPriceBreak {
	var result []types.UnifiedPriceBreak

	for _, pb := range priceBreaks {
		base := pb.Price
		markup := base * 1.10
		targetPurch := base * 0.82
		costDelivery := targetPurch + deliveryCoef
		targetSales := costDelivery + markup

		result = append(result, types.UnifiedPriceBreak{
			Quantity:              pb.Quantity,
			Price:                 utils.Round(base, 2),
			Currency:              currency,
			CostWithDelivery:      utils.Round(costDelivery, 2),
			TargetPricePurchasing: utils.Round(targetPurch, 2),
			TargetPriceSales:      utils.Round(targetSales, 2),
		})
	}

	return result
}

func formatDeliveryTime(days int, defaultValue string) string {
	if days <= 0 {
		return defaultValue
	}

	if days <= 7 {
		return "1-2 недели"
	}

	weeks := float64(days) / 7.0

	minWeeks := int(weeks)
	maxWeeks := minWeeks + 1

	if days%7 == 0 {
		return fmt.Sprintf("%d недели", minWeeks)
	}

	return fmt.Sprintf("%d-%d недели", minWeeks, maxWeeks)
}
