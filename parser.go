package api

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

//Cube struct
type Cube struct {
	XMLName  xml.Name `xml:"Cube"`
	Currency string   `xml:"currency,attr"`
	Rate     float32  `xml:"rate,attr"`
}

//Cubes struct
type Cubes struct {
	XMLName xml.Name `xml:"Cube"`
	Time    string   `xml:"time,attr"`
	Rates   []Cube   `xml:"Cube"`
}

//ParentCube Struct
type ParentCube struct {
	XMLName xml.Name `xml:"Cube"`
	Cubes   []Cubes  `xml:"Cube"`
}

//Envelope struct
type Envelope struct {
	XMLName xml.Name   `xml:"Envelope"`
	Cube    ParentCube `xml:"Cube"`
}

func downloadXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

//ParseXML to string
func ParseXML(url string) ([]Cubes, error) {
	b, err := downloadXML(URL)
	if err != nil {
		return nil, err
	}

	e := Envelope{}

	if err := xml.Unmarshal(b, &e); err != nil {
		return nil, err
	}

	return e.Cube.Cubes, nil
}
