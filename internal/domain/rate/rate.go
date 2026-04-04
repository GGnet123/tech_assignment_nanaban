package rate

type Result struct {
	Bid       float64
	Ask       float64
	Timestamp int64
}

type RateSide string

var SideBid RateSide = "bid"
var SideAsk RateSide = "ask"

type SaveRate struct {
	Price     float64
	Side      RateSide
	Timestamp int64
}
