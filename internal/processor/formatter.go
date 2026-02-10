package processor

import (
	"strconv"

	"dynamic-pricing-tool-ru/internal/types"
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

// FormatGetchipsData преобразует сырые данные Getchips в упрощенный формат
func FormatGetchipsData(raw *types.GetchipsResponse, partNumber string) *types.SimplifiedGetchipsData {
	if raw == nil || len(raw.Data) == 0 {
		return nil
	}

	data := raw.Data[0]

	// Определяем валюту
	currency := "USD"
	if data.Currency == 1 {
		currency = "USD"
	} else if data.Currency == 2 {
		currency = "EUR"
	}

	// Преобразуем PriceBreak
	priceBreaks := make([]types.PriceBreak, 0, len(data.PriceBreak))
	for _, pb := range data.PriceBreak {
		priceBreaks = append(priceBreaks, types.PriceBreak{
			Quantity: pb.Quantity,
			Price:    pb.Price,
		})
	}

	return &types.SimplifiedGetchipsData{
		PartNumber:   data.Title,
		AvailableQty: data.Quantity,
		Brand:        data.Brand,
		Supplier:     data.Donor,
		Price:        data.Price,
		Currency:     currency,
		LeadTimeDays: data.Orderdays,
		PriceBreaks:  priceBreaks,
		Packaging:    data.Packaging,
		MinOrderQty:  data.Minq,
	}
}

// FormatEfindData преобразует сырые данные Efind в упрощенный формат
func FormatEfindData(raw *types.EfindResponse) *types.SimplifiedEfindData {
	if raw == nil || len(*raw) == 0 {
		return nil
	}

	// Берем первый результат
	result := (*raw)[0]

	if len(result.Rows) == 0 {
		return nil
	}

	row := result.Rows[0]

	// Парсим количество
	availableQty := toInt(row.Stock)
	if availableQty < 0 {
		availableQty = 0
	}

	inStock := availableQty > 0

	// Парсим минимальное количество заказа
	minOrderQty := toInt(row.Moq)

	// Извлекаем цену (первая цена из массива)
	var price float64
	if len(row.Price) > 0 && len(row.Price[0]) > 2 {
		if priceVal, ok := row.Price[0][2].(float64); ok {
			price = priceVal
		}
	}

	// Преобразуем прайс-брейки
	priceBreaks := make([]types.PriceBreak, 0, len(row.Price))
	for _, p := range row.Price {
		if len(p) >= 3 {
			qty := toInt(p[0])

			price, ok := p[2].(float64)
			if !ok {
				continue
			}

			priceBreaks = append(priceBreaks, types.PriceBreak{
				Quantity: qty,
				Price:    price,
			})
		}
	}

	supplierInfo := types.SupplierInfo{
		Name:     result.StockData.Title,
		City:     result.StockData.City,
		Country:  result.StockData.Country,
		Email:    result.StockData.ContactEmail,
		Phones:   result.StockData.ContactPhones,
		Website:  result.StockData.Site,
		MinOrder: result.StockData.MinOrder,
	}

	return &types.SimplifiedEfindData{
		PartNumber:   row.Part,
		Supplier:     result.StockData.Title,
		AvailableQty: availableQty,
		Price:        price,
		Currency:     row.Cur,
		PriceBreaks:  priceBreaks,
		SupplierInfo: supplierInfo,
		InStock:      inStock,
		MinOrderQty:  minOrderQty,
	}
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

func FormatPromelecData(data types.PromelecResponse) []types.SimplifiedPromelecData {
	var result []types.SimplifiedPromelecData

	for _, item := range data {
		s := types.SimplifiedPromelecData{
			PartNumber:   item.Name,
			Manufacturer: item.ProducerName,
			Package:      item.Package,
			Stock:        item.Quant,
			MOQ:          item.Moq,
		}

		for _, pb := range item.Pricebreaks {
			s.Prices = append(s.Prices, types.PriceBreak{
				Quantity: pb.Quant,
				Price:    pb.Price,
			})
		}

		result = append(result, s)
	}

	return result
}
