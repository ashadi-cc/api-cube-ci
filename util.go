package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func dbConnect() (*sql.DB, error) {
	strCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		MYSQL_USERNAME,
		MYSQL_PASSWORD,
		MYSQL_HOST,
		MYSQL_PORT,
		MYSQL_DB,
	)

	db, err := sql.Open("mysql", strCon)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func bulkCubeInsert(cubes []Cubes, db *sql.DB) error {
	if len(cubes) == 0 {
		return fmt.Errorf("empty cubes %s", "0")
	}
	valueStrings := []string{}
	valueArgs := []interface{}{}

	for _, c := range cubes {
		for _, r := range c.Rates {
			valueStrings = append(valueStrings, "(?,?,?)")
			valueArgs = append(valueArgs, c.Time)
			valueArgs = append(valueArgs, r.Currency)
			valueArgs = append(valueArgs, r.Rate)
		}
	}

	smt := "INSERT INTO cubes (rate_date,currency,rate) values %s"
	smt += " ON DUPLICATE KEY UPDATE rate=VALUES(rate)"
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))
	_, err := db.Exec(smt, valueArgs...)

	return err
}

//RespondJSON send response as json
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

//RespondError send error response as json
func RespondError(w http.ResponseWriter, code int, message ...*errorData) {
	RespondJSON(w, code, map[string][]*errorData{"error": message})
}
