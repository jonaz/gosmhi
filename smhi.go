package gosmhi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

const URL = "http://opendata-download-metfcst.smhi.se/api/category/pmp1.5g/version/1/geopoint/lat/{{.Latitude}}/lon/{{.Longitude}}/data.json"

type timeSerie struct {
	validTime string
	t         float64
	msl       float64
	vis       float64
	wd        int
	ws        float64
	r         float64
	tstm      float64
	tcc       float64
	lcc       float64
	mcc       float64
	hcc       float64
	gust      float64
	pis       float64
	pit       float64
	pcat      float64
}

type response struct {
	Lat           float64     `json: "lat"`
	Lon           float64     `json: "lon"`
	ReferenceTime string      `json: "referenceTime"`
	Timeseries    []timeSerie `json: "timeSeries"`
}

type smhi struct {
	request *request
}

type request struct {
	Latitude  string
	Longitude string
	url       string
}

func (req *request) getUrl(latitude string, longitude string) string { // {{{
	req.Latitude = latitude
	req.Longitude = longitude

	tmpl, err := template.New("Url").Parse(URL)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	//err = tmpl.Execute(os.Stdout, req)
	err = tmpl.Execute(&doc, req)
	if err != nil {
		panic(err)
	}
	return doc.String()
} // }}}

func New() *smhi { // {{{
	request := &request{}
	return &smhi{request}
} // }}}

func (smhi *smhi) GetByLatLong(latitude string, longitude string) {
	//check if we can get from cache?
	smhi.request.url = smhi.request.getUrl(latitude, longitude)

	fmt.Println(smhi.doRequest())
}

func (smhi *smhi) doRequest() (*response, error) {

	resp, err := http.Get(smhi.request.url)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	response := &response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return response, nil
}
