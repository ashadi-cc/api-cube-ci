package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//App wrapper
type App struct {
	DB      *sql.DB
	Handler *http.ServeMux
}

//ImportXML import xml from url
func (app *App) ImportXML() (int, error) {
	cubes, err := ParseXML(URL)
	if err != nil {
		return 0, fmt.Errorf("Could not download xml from %s, got error: %s", URL, err.Error())
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

func (app *App) getRateByDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		RespondError(w, http.StatusBadRequest, newAPIError(errMethodNotAllowed))
		return
	}

	rdate := strings.TrimPrefix(r.URL.Path, "/rates/")
	if rdate == "" {
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

func (app *App) setRouter() {
	app.Handler = http.NewServeMux()
	app.Handler.HandleFunc("/rates/latest", app.getLatestRate)
	app.Handler.HandleFunc("/rates/", app.getRateByDate)
	app.Handler.HandleFunc("/rates/analyze", app.getAnalizeReport)
}

//Init init database
func (app *App) Init() {
	db, err := dbConnect()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = db
	//start import
	log.Println("Start importing rates from ", URL)
	go app.ImportXML()
	log.Println("Set router")
	app.setRouter()
}

//Run web server
func (app *App) Run() {
	log.Printf("App listening on port %s \n", APP_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", APP_PORT), app.Handler))
}
