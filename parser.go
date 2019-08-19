package api

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

//Cube represent detail of Cubes node
type Cube struct {
	XMLName  xml.Name `xml:"Cube"`
	Currency string   `xml:"currency,attr"`
	Rate     float32  `xml:"rate,attr"`
}

//Cubes represent detail of Parentcube node
type Cubes struct {
	XMLName xml.Name `xml:"Cube"`
	Time    string   `xml:"time,attr"`
	Rates   []Cube   `xml:"Cube"`
}

//ParentCube represent main node of Cube
type ParentCube struct {
	XMLName xml.Name `xml:"Cube"`
	Cubes   []Cubes  `xml:"Cube"`
}

//Envelope represent root XML Node
type Envelope struct {
	XMLName xml.Name   `xml:"Envelope"`
	Cube    ParentCube `xml:"Cube"`
}

//downloadXML download xml from given url
func downloadXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

//ParseXML parse xml to cubes model
func ParseXML(b []byte) ([]Cubes, error) {

	e := Envelope{}

	if err := xml.Unmarshal(b, &e); err != nil {
		return nil, err
	}

	return e.Cube.Cubes, nil
}
