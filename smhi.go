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
	ValidTime string  `json: "validTime"` //Time
	T         float64 `json: "t"`         //Temperature celcius
	Msl       float64 `json: "msl"`       //Pressure reduced to MSL hPa
	Vis       float64 `json: "vis"`       //Visibility km
	Wd        int     `json: "wd"`        //wind direction degrees
	Ws        float64 `json: "ws"`        //wind velocity m/s
	R         int     `json: "r"`         //Relative humidity %
	Tstm      int     `json: "tstm"`      //Probability thunderstorm %
	Tcc       int     `json: "tcc"`       //Total cloud cover 0-8
	Lcc       int     `json: "lcc"`       //Low cloud cover 0-8
	Mcc       int     `json: "mcc"`       //Medium cloud cover 0-8
	Hcc       int     `json: "hcc"`       //high cloud cover 0-8
	Gust      float64 `json: "gust"`      //Wind gust m/s
	Pis       float64 `json: "pis"`       //Precipitation intensity snow mm/h
	Pit       float64 `json: "pit"`       //Precipitation intensity total mm/h
	Pcat      int     `json: "pcat"`      //Category of precipitation, 0 no, 1 snow, 2 snow and rain, 3 rain, 4 drizzle, 5, freezing rain, 6 freezing drizzle
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
