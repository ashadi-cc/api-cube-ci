package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	totalCubes int
	DB         *sql.DB
	app        *App
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Handler.ServeHTTP(rr, req)
	return rr
}

func TestDownloadXML(t *testing.T) {
	_, err := downloadXML(URL)
	if err != nil {
		t.Fatalf("it should download xml from %s but got error %s", URL, err.Error())
	}
}

func TestParseXML(t *testing.T) {
	cubes, err := ParseXML(URL)
	if err != nil {
		t.Fatalf("Error message it should nil but got %s", err.Error())
	}
	if cubes == nil {
		t.Fatalf("It should get cubes list but got empty result")
	}

	cube := cubes[0]
	if cube.Time == "" {
		t.Fatalf("the time its should not empty")
	}

	if len(cube.Rates) == 0 {
		t.Fatalf("rates its should not empty record")
	}

	totalCubes = len(cubes)
}

func TestGetEnv(t *testing.T) {
	val := getEnv("USER_ENV", "root")
	if val != "root" {
		t.Fatalf("User_ENV variable its should %s but got %s", "root", val)
	}
}

func TestConnDatabase(t *testing.T) {
	db, err := dbConnect()
	if err != nil {
		t.Fatalf("its should connected to mysql but got error : %s", err.Error())
	}
	DB = db
}

func TestBulkInsertCubes(t *testing.T) {

	if _, err := DB.Exec("DELETE FROM cubes"); err != nil {
		t.Fatalf("it should not error when erase cubes table %s", err.Error())
	}
	rates := []Cube{Cube{Currency: "USD", Rate: 10.23}}
	cubes := []Cubes{Cubes{Time: "2019-02-03", Rates: rates}}

	err := bulkCubeInsert(cubes, DB)

	if err != nil {
		t.Fatalf("it should not raise error when insert cube but got %s", err.Error())
	}
}

func runImport(db *sql.DB, t *testing.T) {
	app := &App{DB: DB}

	total, err := app.ImportXML()

	if err != nil {
		t.Fatalf("it should not error when import xml from %s, but got error: %s", URL, err.Error())
	}

	if total != totalCubes {
		t.Fatal("it should not return zero record when import xml")
	}
}

func TestImportXML(t *testing.T) {

	if _, err := DB.Exec("DELETE FROM cubes"); err != nil {
		t.Fatalf("it should not error when erase cubes table %s", err.Error())
	}

	runImport(DB, t)

}

func TestDuplicateImport(t *testing.T) {
	runImport(DB, t)
}

func TestGetLatestResponse(t *testing.T) {
	response, err := GetLatestRate(DB)

	if err != nil {
		t.Fatalf("it should not get error when get latest rate : %s", err.Error())
	}

	if response.Base != "EUR" {
		t.Fatalf("base response should EUR but got %s", response.Base)
	}

	if len(response.Rates) == 0 {
		t.Fatal("its should not get empty record for rates")
	}

}

func TestGetAnalizeRport(t *testing.T) {
	response, err := GetAnalyzeRate(DB)

	if err != nil {
		t.Fatalf("it should not get error when get latest rate : %s", err.Error())
	}

	if response.Base != "EUR" {
		t.Fatalf("base response should EUR but got %s", response.Base)
	}

	if len(response.RateAnalyze) == 0 {
		t.Fatal("its should not get empty record for rates")
	}
}

func TestInitWebServer(t *testing.T) {
	app = &App{}
	app.Init()
}

func TestGetLatestRate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/rates/latest", nil)
	response := executeRequest(req)
	if response.Code != http.StatusOK {
		t.Fatalf("it should get status 200 but got %d", response.Code)
	}
}

func TestGetAnalizeRate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/rates/analyze", nil)
	response := executeRequest(req)
	if response.Code != http.StatusOK {
		t.Fatalf("it should get status 200 but got %d", response.Code)
	}
}

func TestGetRateByDate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/rates/2019-05-20", nil)
	response := executeRequest(req)
	if response.Code != http.StatusOK {
		t.Fatalf("it should get status 200 but got %d", response.Code)
	}
}
