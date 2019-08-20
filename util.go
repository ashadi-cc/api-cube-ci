package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

//dbConnect Connect to Mysql database
func dbConnect(config *DbConfig) (*sql.DB, error) {
	strCon := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", strCon)

	if err != nil {
		return nil, err
	}

	return db, nil
}

// bulkCubeInsert store exchanges rates to database as bulk insert to speed up the process
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

//RespondJSON send response as json format
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errinternal := []*errorData{newAPIError(errInternalServerError)}
		errMsg := map[string][]*errorData{"errors": errinternal}
		r, _ := json.Marshal(errMsg)
		w.Write([]byte(r))
		return
	}

	w.WriteHeader(status)
	w.Write([]byte(response))
}

//RespondError send error response as json
func RespondError(w http.ResponseWriter, code int, message ...*errorData) {
	RespondJSON(w, code, map[string][]*errorData{"errors": message})
}

//validating input should have format YYYY-MM-DD
func checkValidDate(date string) bool {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	return re.MatchString(date)
}
