package models

//todo: add more fields so ticker event can be also serialized properly
type Event struct {
	Event              string  `json:"e"`
	EventTime          int64   `json:"E"`
	Pair               string  `json:"s"`
	TradeId            int64   `json:"t"`
	TradePrice         float32 `json:"p,string"`
	TradeVolume        float32 `json:"q,string"`
	ByOrderId          int64   `json:"b"`
	SellOrderId        int64   `json:"a"`
	TradeCompletedTime int64   `json:"T"`
}
