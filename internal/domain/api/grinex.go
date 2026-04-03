package api

type OrderBook struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []Order `json:"asks"`
	Bids      []Order `json:"bids"`
}

type Order struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
	Volume string `json:"volume"`
}
