package types

type Request struct {
	Mapping map[string]string `json:"mapping"`
	Data    [][]string        `json:"data"`
	Mode    string            `json:"mode"`
}

type PartData struct {
	PartNumber string
	Quantity   string
	RowIndex   int
}

type GetchipsResponse struct {
	Data []struct {
		Title         string       `json:"title"`
		Quantity      int          `json:"quantity"`
		SPack         int          `json:"sPack"`
		DonorID       interface{}  `json:"donorID"`
		Donor         string       `json:"donor"`
		Folddivision  int          `json:"folddivision"`
		Minq          int          `json:"minq"`
		Brand         string       `json:"brand"`
		Orderdays     int          `json:"orderdays"`
		Price         float64      `json:"price"`
		EQuantity     int          `json:"eQuantity"`
		SearchWord    string       `json:"search_word"`
		Currency      int          `json:"currency"`
		Match         int          `json:"match"`
		PriceBreak    []PriceBreak `json:"priceBreak"`
		QuantityPrice float64      `json:"quantityPrice"`
		Packaging     string       `json:"packaging"`
	} `json:"data"`
}

type EfindResponse []struct {
	Filial interface{} `json:"filial"`
	Finish float64     `json:"finish"`
	Rows   []struct {
		Part    string          `json:"part"`
		Cur     string          `json:"cur"`
		Instock bool            `json:"instock"`
		Price   [][]interface{} `json:"price"`

		Moq   interface{} `json:"moq"`
		Mpq   interface{} `json:"mpq"`
		Stock interface{} `json:"stock"`
		Od    interface{} `json:"od"`
	} `json:"rows"`

	StockID   int `json:"stock_id"`
	StockData struct {
		City          string   `json:"city"`
		ContactEmail  string   `json:"contact_email"`
		ContactPhones []string `json:"contact_phones"`
		Country       string   `json:"country"`
		MinOrder      string   `json:"min_order"`
		RegionID      string   `json:"region_id"`
		Site          string   `json:"site"`
		Title         string   `json:"title"`
		TitleEn       string   `json:"title_en"`
	} `json:"stockdata"`
}

// Упрощенная структура для хранения основных данных
type SimplifiedGetchipsData struct {
	PartNumber   string       `json:"partNumber"`
	AvailableQty int          `json:"availableQty"`
	Brand        string       `json:"brand"`
	Supplier     string       `json:"supplier"`
	Price        float64      `json:"price"`
	Currency     string       `json:"currency"`
	LeadTimeDays int          `json:"leadTimeDays"`
	PriceBreaks  []PriceBreak `json:"priceBreaks"`
	Packaging    string       `json:"packaging"`
	MinOrderQty  int          `json:"minOrderQty"`
}

type SimplifiedEfindData struct {
	PartNumber   string       `json:"partNumber"`
	Supplier     string       `json:"supplier"`
	AvailableQty int          `json:"availableQty"`
	Price        float64      `json:"price"`
	Currency     string       `json:"currency"`
	PriceBreaks  []PriceBreak `json:"priceBreaks"`
	SupplierInfo SupplierInfo `json:"supplierInfo"`
	InStock      bool         `json:"inStock"`
	MinOrderQty  int          `json:"minOrderQty"`
}

type PriceBreak struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total,omitempty"`
}

type SupplierInfo struct {
	Name     string   `json:"name"`
	City     string   `json:"city"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Phones   []string `json:"phones"`
	Website  string   `json:"website"`
	MinOrder string   `json:"minOrder"`
}

type CombinedResult struct {
	PartNumber    string                   `json:"partNumber"`
	RequestedQty  int                      `json:"requestedQty"`
	RowIndex      int                      `json:"rowIndex"`
	Getchips      *SimplifiedGetchipsData  `json:"getchips,omitempty"`
	Efind         *SimplifiedEfindData     `json:"efind,omitempty"`
	Promelec      []SimplifiedPromelecData `json:"promelec,omitempty"`
	GetchipsRaw   *GetchipsResponse        `json:"getchipsRaw,omitempty"`
	EfindRaw      *EfindResponse           `json:"efindRaw,omitempty"`
	GetchipsError string                   `json:"getchipsError,omitempty"`
	EfindError    string                   `json:"efindError,omitempty"`
	PromelecError string                   `json:"promelecError,omitempty"`
	Timestamp     string                   `json:"timestamp"`
}

type APIResponse struct {
	PartNumber   string
	Quantity     string
	RequestedQty int
	GetchipsData *GetchipsResponse
	EfindData    *EfindResponse
	PromelecData PromelecResponse
	GetchipsErr  error
	EfindErr     error
	PromelecErr  error
	RowIndex     int
}

type AnalysisResult struct {
	TotalParts       int               `json:"totalParts"`
	GetchipsSuccess  int               `json:"getchipsSuccess"`
	GetchipsInStock  int               `json:"getchipsInStock"`
	EfindSuccess     int               `json:"efindSuccess"`
	EfindInStock     int               `json:"efindInStock"`
	BothAPIsSuccess  int               `json:"bothAPIsSuccess"`
	PriceComparisons []PriceComparison `json:"priceComparisons,omitempty"`
}

type PriceComparison struct {
	PartNumber        string  `json:"partNumber"`
	GetchipsPrice     float64 `json:"getchipsPrice"`
	EfindPrice        float64 `json:"efindPrice"`
	BetterPrice       string  `json:"betterPrice"`
	DifferencePercent float64 `json:"differencePercent"`
}

// ================= PROMELEC =================

type PromelecResponse []PromelecItem

type PromelecItem struct {
	ItemID       int    `json:"item_id"`
	Name         string `json:"name"`
	ProducerName string `json:"producer_name"`
	Package      string `json:"package"`
	Quant        int    `json:"quant"`
	Moq          int    `json:"moq"`
	Munit        string `json:"munit"`
	Pricebreaks  []struct {
		Quant int     `json:"quant"`
		Price float64 `json:"price"`
	} `json:"pricebreaks"`
}

type SimplifiedPromelecData struct {
	PartNumber   string
	Manufacturer string
	Package      string
	Stock        int
	MOQ          int
	Prices       []PriceBreak
}
