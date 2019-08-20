package api

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	totalCubes int
	DB         *sql.DB
	app        *App
	config     *Config
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Handler.ServeHTTP(rr, req)
	return rr
}

func TestMain(m *testing.M) {
	config = LoadConfig()
	code := m.Run()
	os.Exit(code)
}

func TestSplitCubes(t *testing.T) {

	if c := splitCubes(make([]Cubes, 0), 10); len(c) != 0 {
		t.Fatalf("total cubes should get 0 but got %d", len(c))
	}
	if c := splitCubes(make([]Cubes, 3), 10); len(c) != 1 {
		t.Fatalf("total cubes should get 1 but got %d", len(c))
	}

	if c := splitCubes(make([]Cubes, 10), 10); len(c) != 1 {
		t.Fatalf("total cubes should get 1 but got %d", len(c))
	}

	if c := splitCubes(make([]Cubes, 20), 10); len(c) != 2 {
		t.Fatalf("total cubes should get 2 but got %d", len(c))
	}

	if c := splitCubes(make([]Cubes, 11), 10); len(c) != 2 {
		t.Fatalf("total cubes should get 2 but got %d", len(c))
	}

	c := splitCubes(make([]Cubes, 13), 10)

	if l := len(c[0]); l != 10 {
		t.Fatalf("total cubes[0], should 10 but got %d", l)
	}

	if l := len(c[1]); l != 3 {
		t.Fatalf("total cubes[0], should 3 but got %d", l)
	}
}

func TestValidDateFunction(t *testing.T) {
	date := "2019-01-02"
	if v := checkValidDate(date); !v {
		t.Fatalf("%s it should valid date", date)
	}

	date = "02-01-02"
	if v := checkValidDate(date); v {
		t.Fatalf("%s it should not valid date", date)
	}

	date = "yyyy-mm-dd"
	if v := checkValidDate(date); v {
		t.Fatalf("%s it should not valid date", date)
	}
}

func TestDownloadXML(t *testing.T) {
	_, err := downloadXML(config.App.XMLUrl)
	if err != nil {
		t.Fatalf("it should download xml from %s but got error %s", config.App.XMLUrl, err.Error())
	}
}

func TestParseXML(t *testing.T) {
	b, _ := downloadXML(config.App.XMLUrl)
	cubes, err := ParseXML(b)
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
	db, err := dbConnect(config.Db)
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
	app := &App{DB: DB, Config: config}
	_, err := app.ImportXML()
	if err != nil {
		t.Fatalf("it should not error when import xml from %s, but got error: %s", config.App.XMLUrl, err.Error())
	}
}

func TestImportXML(t *testing.T) {
	if _, err := DB.Exec("DELETE FROM cubes"); err != nil {
		t.Fatalf("it should not error when erase cubes table %s", err.Error())
	}

	runImport(DB, t)
}

func TestDuplicateImport(t *testing.T) {
	DB.QueryRow("SELECT COUNT(id) as t FROM cubes").Scan(&totalCubes)
	runImport(DB, t)

	var total int
	DB.QueryRow("SELECT COUNT(id) as t FROM cubes").Scan(&total)
	if total != totalCubes {
		t.Fatalf("total its should %d but got %d", totalCubes, total)
	}

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
	app = New()
	app.Init()
	//set test configuration
	app.Config = config
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

func TestGetRateByInvalidDate(t *testing.T) {
	req, _ := http.NewRequest("GET", "/rates/YYYY-MM-DD", nil)
	response := executeRequest(req)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("it should get status 400 but got %d", response.Code)
	}
}
