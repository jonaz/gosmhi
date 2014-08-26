package gosmhi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	//"strings"
	"text/template"
	"time"
)

const URL = "http://opendata-download-metfcst.smhi.se/api/category/pmp1.5g/version/1/geopoint/lat/{{.Latitude}}/lon/{{.Longitude}}/data.json"

type timeSerie struct {
	ValidTime string    `json: "validTime"` //Time
	Time      time.Time //Time in go
	T         float64   `json: "t"`    //Temperature celcius
	Msl       float64   `json: "msl"`  //Pressure reduced to MSL hPa
	Vis       float64   `json: "vis"`  //Visibility km
	Wd        int       `json: "wd"`   //wind direction degrees
	Ws        float64   `json: "ws"`   //wind velocity m/s
	R         int       `json: "r"`    //Relative humidity %
	Tstm      int       `json: "tstm"` //Probability thunderstorm %
	Tcc       int       `json: "tcc"`  //Total cloud cover 0-8
	Lcc       int       `json: "lcc"`  //Low cloud cover 0-8
	Mcc       int       `json: "mcc"`  //Medium cloud cover 0-8
	Hcc       int       `json: "hcc"`  //high cloud cover 0-8
	Gust      float64   `json: "gust"` //Wind gust m/s
	Pis       float64   `json: "pis"`  //Precipitation intensity snow mm/h
	Pit       float64   `json: "pit"`  //Precipitation intensity total mm/h
	Pcat      int       `json: "pcat"` //Category of precipitation, 0 no, 1 snow, 2 snow and rain, 3 rain, 4 drizzle, 5, freezing rain, 6 freezing drizzle
}

type response struct {
	Lat           float64     `json: "lat"`
	Lon           float64     `json: "lon"`
	ReferenceTime string      `json: "referenceTime"`
	TimeSeries    []timeSerie `json: "timeSeries"`
}

//time,_:=time.Parse(time.RFC3339, "2014-08-26T07:00:00Z")
//fmt.Println(time)

func (resp *response) GetMaxTempByDate(date time.Time) (float64, error) {
	if resp == nil || resp.TimeSeries == nil {
		return 0, errors.New("Invalid response")
	}
	temp := resp.TimeSeries[0].T
	for _, row := range resp.TimeSeries {
		//if strings.Contains(row.ValidTime, date) {
		if date.Year() == row.Time.Year() && date.Day() == row.Time.Day() && date.Month() == row.Time.Month() {
			if row.T > temp {
				temp = row.T
			}
		}
	}
	return temp, nil
}
func (resp *response) GetMinTempByDate(date time.Time) (float64, error) {
	if resp == nil || resp.TimeSeries == nil {
		return 0, errors.New("Invalid response")
	}
	temp := resp.TimeSeries[0].T
	for _, row := range resp.TimeSeries {
		//if strings.Contains(row.ValidTime, date) {
		if date.Year() == row.Time.Year() && date.Day() == row.Time.Day() && date.Month() == row.Time.Month() {
			if row.T < temp {
				temp = row.T
			}
		}
	}
	return temp, nil
}

type smhi struct {
	request  *request
	response *response
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
	return &smhi{request, nil}
} // }}}

func (smhi *smhi) GetByLatLong(latitude string, longitude string) *response {
	//check if we can get from cache?
	smhi.request.url = smhi.request.getUrl(latitude, longitude)

	resp, err := smhi.doRequest()
	if err != nil {
		fmt.Println(err)
	}

	smhi.response = resp
	return resp
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

	parseValidTime(response)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return response, nil
}

func parseValidTime(resp *response) {
	var t time.Time
	var err error
	for key, row := range resp.TimeSeries {
		t, err = time.Parse(time.RFC3339, row.ValidTime)
		if err != nil {
			fmt.Println("error parsing time")
		}
		resp.TimeSeries[key].Time = t
	}
}
