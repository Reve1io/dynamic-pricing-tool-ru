package types

type Request struct {
	Mapping map[string]string `json:"mapping"`
	Data    [][]string        `json:"data"`
	Mode    string            `json:"mode"`
}

type PartData struct {
	PartNumber string
	Quantity   string
}

type GetchipsResponse struct {
	Data []struct {
		Title      string  `json:"title"`
		Quantity   int     `json:"quantity"`
		Brand      string  `json:"brand"`
		Price      float64 `json:"price"`
		Currency   int     `json:"currency"`
		Donor      string  `json:"donor"`
		OrderDays  int     `json:"orderdays"`
		PriceBreak []struct {
			Quantity int     `json:"quantity"`
			Price    float64 `json:"price"`
			Summ     float64 `json:"summ"`
		} `json:"priceBreak"`
		Packaging string `json:"packaging"`
	} `json:"data"`
}

type APIResponse struct {
	PartNumber string
	Data       interface{}
	Error      error
}

type Config struct {
	ServerPort    string
	GetchipsURL   string
	GetchipsToken string
	RedisAddr     string
	RabbitMQURL   string
	ChunkSize     int
}
