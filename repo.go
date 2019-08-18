package main

import (
	"database/sql"
)

//GetRateByDate get rate by date
func GetRateByDate(rdate string, db *sql.DB) (LatestResponse, error) {
	response := LatestResponse{}
	rows, err1 := db.Query("SELECT currency,rate FROM cubes WHERE rate_date = ? ORDER BY currency", rdate)
	if err1 != nil {
		return response, err1
	}

	rates := Rate{}
	for rows.Next() {
		cube := Cube{}
		if err := rows.Scan(&cube.Currency, &cube.Rate); err != nil {
			return response, err
		}
		rates[cube.Currency] = cube.Rate
	}

	response.Base, response.Rates = "EUR", rates

	return response, nil
}

//GetLatestRate query
func GetLatestRate(db *sql.DB) (LatestResponse, error) {
	var rdate string
	err := db.QueryRow("SELECT rate_date FROM cubes ORDER BY rate_date DESC LIMIT 0,1").Scan(&rdate)
	if err != nil {
		return LatestResponse{}, err
	}

	return GetRateByDate(rdate, db)
}

//GetAnalyzeRate get analyze report of rate
func GetAnalyzeRate(db *sql.DB) (AnalizeResponse, error) {
	response := AnalizeResponse{}
	sql := "SELECT MIN(rate) AS i, MAX(rate) AS a, AVG(rate) AS v, currency FROM cubes GROUP BY currency"

	rows, err1 := db.Query(sql)
	if err1 != nil {
		return response, err1
	}

	c, a := "", CurrencyResponse{}

	for rows.Next() {
		r := Summary{}
		if err := rows.Scan(&r.Min, &r.Max, &r.Avg, &c); err != nil {
			return response, err
		}
		a[c] = r
	}
	response.Base = "EUR"
	response.RateAnalyze = a
	return response, nil
}
