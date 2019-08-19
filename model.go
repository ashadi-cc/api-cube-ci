package api

//Rate represent detail exchange rate model
type Rate map[string]float32

//LatestResponse represent exchange rate response
type LatestResponse struct {
	Base  string `json:"base"`
	Rates Rate   `json:"rates"`
}

//Summary represent summary response
type Summary struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
	Avg float32 `json:"avg"`
}

//CurrencyResponse represent currency response
type CurrencyResponse map[string]Summary

//AnalizeResponse represent analyze response
type AnalizeResponse struct {
	Base        string           `json:"base"`
	RateAnalyze CurrencyResponse `json:"rates_analyze"`
}
