package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//App Api wrapper
type App struct {
	DB      *sql.DB
	Handler *http.ServeMux
	Config  *Config
}

//ImportResult struct
type ImportResult struct {
	Total int
	Error error
}

//ImportChan  channel
type ImportChan chan ImportResult

//ImportXML Download XML and store to database
func (app *App) ImportXML(ch ImportChan) {
	chanResult := ImportResult{}
	//downloax xml
	b, err := downloadXML(app.Config.App.XMLUrl)
	if err != nil {
		chanResult.Total, chanResult.Error = 0, fmt.Errorf("Could not download xml from %s, got error: %s", app.Config.App.XMLUrl, err.Error())
		ch <- chanResult
		return
	}

	cubes, err := ParseXML(b)
	if err != nil {
		chanResult.Total, chanResult.Error = 0, fmt.Errorf("Could not parse xml, got error: %s", err.Error())
		ch <- chanResult
		return
	}

	splitCubes, totalRecord := splitCubes(cubes, 500), 0

	for _, c := range splitCubes {
		if err := bulkCubeInsert(c, app.DB); err != nil {
			log.Fatalf("Error when insert cubes %s", err.Error())
		}
		totalRecord += len(c)
	}

	log.Println("import xml done!")
	chanResult.Total, chanResult.Error = totalRecord, nil
	ch <- chanResult
}

//getLatestRate /rates/latest endpoint implementation
func (app *App) getLatestRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		RespondError(w, http.StatusBadRequest, newAPIError(errMethodNotAllowed))
		return
	}

	payload, err := GetLatestRate(app.DB)

	if err != nil {
		RespondError(w, http.StatusInternalServerError, newAPIError(errInternalServerError))
		return
	}

	RespondJSON(w, http.StatusOK, payload)
}

//getAnalizeReport /rates/analyze endpoint implementation
func (app *App) getAnalizeReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		RespondError(w, http.StatusBadRequest, newAPIError(errMethodNotAllowed))
		return
	}

	payload, err := GetAnalyzeRate(app.DB)

	if err != nil {
		RespondError(w, http.StatusInternalServerError, newAPIError(errInternalServerError))
		return
	}

	RespondJSON(w, http.StatusOK, payload)
}

//getRateByDate /rates/YYYY-MM-DD endpoint implementation
func (app *App) getRateByDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		RespondError(w, http.StatusBadRequest, newAPIError(errMethodNotAllowed))
		return
	}

	rdate := strings.TrimPrefix(r.URL.Path, "/rates/")
	if !checkValidDate(rdate) {
		RespondError(w, http.StatusBadRequest, newAPIError(errBadRequest))
		return
	}

	payload, err := GetRateByDate(rdate, app.DB)

	if err != nil {
		log.Printf("get rates by date %s got error: %s \n", rdate, err.Error())
		RespondError(w, http.StatusInternalServerError, newAPIError(errInternalServerError))
		return
	}

	RespondJSON(w, http.StatusOK, payload)
}

//setRouter Configure endpoint routes
func (app *App) setRouter() {
	app.Handler = http.NewServeMux()
	app.Handler.HandleFunc("/rates/latest", app.getLatestRate)
	app.Handler.HandleFunc("/rates/", app.getRateByDate)
	app.Handler.HandleFunc("/rates/analyze", app.getAnalizeReport)
}

//attach database connection
func (app *App) setDatabase() {
	db, err := dbConnect(app.Config.Db)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = db
}

//Init run initial process
func (app *App) Init() {
	//load configuration
	app.Config = LoadConfig()

	//set database
	app.setDatabase()

	//start download and import XML on background
	log.Println("Start importing rates from ", app.Config.App.XMLUrl)
	ch := make(ImportChan)
	go app.ImportXML(ch)
	<-ch
	//Configure endpoint routes
	app.setRouter()
}

//Run http web server
func (app *App) Run() {
	log.Printf("App listening on port %s \n", app.Config.App.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", app.Config.App.Port), app.Handler))
}

//New create new App
func New() *App {
	api := &App{}
	return api
}
