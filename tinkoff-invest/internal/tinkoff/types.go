package tinkoff

type PortfolioPosition struct {
	Figi           string  `json:"figi"`
	Ticker         string  `json:"ticker"`
	Isin           string  `json:"isin"`
	InstrumentType string  `json:"instrumentType"`
	Balance        float64 `json:"balance"`
	Blocked        float64 `json:"blocked"`
}

type PortfolioPayload struct {
	Positions []PortfolioPosition `json:"positions"`
}

type Portfolio struct {
	TrackingId string           `json:"trackingId"`
	Status     string           `json:"status"`
	Payload    PortfolioPayload `json:"payload"`
}
