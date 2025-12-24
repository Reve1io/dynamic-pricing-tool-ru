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
		Title        string  `json:"title"`
		Quantity     int     `json:"quantity"`
		SPrice       float64 `json:"sPrice"`
		DonorID      *string `json:"donorID"`
		Donor        string  `json:"donor"`
		Folddivision int     `json:"folddivision"`
		Minq         int     `json:"minq"`
		Brand        string  `json:"brand"`
		Orderdays    int     `json:"orderdays"`
		Price        float64 `json:"price"`
		EQuantity    int     `json:"eQuantity"`
		SearchWord   string  `json:"search_word"`
		Currency     int     `json:"currency"`
		Match        int     `json:"match"`
		PriceBreak   []struct {
			Quantity int     `json:"quantity"`
			Price    float64 `json:"price"`
			Summ     float64 `json:"summ"`
		} `json:"priceBreak"`
		QuantityPrice float64 `json:"quantityPrice"`
		Packaging     string  `json:"packaging"`
		Manufacturer  string  `json:"manufacturer"`
		Description   string  `json:"description"`
		SPriceRub     float64 `json:"sPriceRub"`
		PriceRub      float64 `json:"priceRub"`
	} `json:"data"`
}

type EfindResponse []struct {
	Filial interface{} `json:"filial"`
	Finish float64     `json:"finish"`
	Rows   []struct {
		Cr      []interface{}   `json:"cr"`
		Cur     string          `json:"cur"`
		Dc      string          `json:"dc"`
		Dlv     string          `json:"dlv"`
		Img     string          `json:"img"`
		Instock bool            `json:"instock"`
		Mfg     string          `json:"mfg"`
		Moq     string          `json:"moq"`
		Mpq     string          `json:"mpq"`
		Note    string          `json:"note"`
		Od      int             `json:"od"`
		Pack    string          `json:"pack"`
		Part    string          `json:"part"`
		Pdf     string          `json:"pdf"`
		Pkg     string          `json:"pkg"`
		Price   [][]interface{} `json:"price"`
		Sku     string          `json:"sku"`
		Stock   string          `json:"stock"`
		Um      string          `json:"um"`
		Url     string          `json:"url"`
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
	PartNumber    string                  `json:"partNumber"`
	RequestedQty  int                     `json:"requestedQty"`
	RowIndex      int                     `json:"rowIndex"`
	Getchips      *SimplifiedGetchipsData `json:"getchips,omitempty"`
	Efind         *SimplifiedEfindData    `json:"efind,omitempty"`
	GetchipsRaw   *GetchipsResponse       `json:"getchipsRaw,omitempty"`
	EfindRaw      *EfindResponse          `json:"efindRaw,omitempty"`
	GetchipsError string                  `json:"getchipsError,omitempty"`
	EfindError    string                  `json:"efindError,omitempty"`
	Timestamp     string                  `json:"timestamp"`
}

type APIResponse struct {
	PartNumber   string
	Quantity     string
	RequestedQty int
	GetchipsData *GetchipsResponse
	EfindData    *EfindResponse
	GetchipsErr  error
	EfindErr     error
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
