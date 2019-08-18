package main

//Rate model
type Rate map[string]float32

//LatestResponse model
type LatestResponse struct {
	Base  string `json:"base"`
	Rates Rate   `json:"rates"`
}

//Summary struct
type Summary struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
	Avg float32 `json:"avg"`
}

//CurrencyResponse model
type CurrencyResponse map[string]Summary

//AnalizeResponse model
type AnalizeResponse struct {
	Base        string           `json:"base"`
	RateAnalyze CurrencyResponse `json:"rates_analyze"`
}
