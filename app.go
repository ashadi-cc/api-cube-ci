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

//ImportXML Download XML and store to database
func (app *App) ImportXML() (int, error) {
	//downloax xml
	b, err := downloadXML(app.Config.App.XMLUrl)
	if err != nil {
		return 0, fmt.Errorf("Could not download xml from %s, got error: %s", app.Config.App.XMLUrl, err.Error())
	}

	cubes, err := ParseXML(b)
	if err != nil {
		return 0, fmt.Errorf("Could not parse xml, got error: %s", err.Error())
	}

	split := 500
	tcubes := []Cubes{}
	totalRecord := 0

	for _, cube := range cubes {
		tcubes = append(tcubes, cube)
		if len(tcubes) >= split {
			if err := bulkCubeInsert(tcubes, app.DB); err != nil {
				log.Fatalf("Error when insert cubes %s", err.Error())
			}
			totalRecord += len(tcubes)
			tcubes = nil
		}
	}

	if err := bulkCubeInsert(tcubes, app.DB); err == nil {
		totalRecord += len(tcubes)
	}

	log.Println("import xml done!")
	return totalRecord, nil
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

//Init run initial process
func (app *App) Init() {
	//load configuration
	app.Config = LoadConfig()

	//attach database connection
	db, err := dbConnect(app.Config.Db)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = db

	//start download and import XML on background
	log.Println("Start importing rates from ", app.Config.App.XMLUrl)
	go app.ImportXML()

	//Configure endpoint routes
	log.Println("Set router")
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
