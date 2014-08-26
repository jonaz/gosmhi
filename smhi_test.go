package gosmhi

import (
	"fmt"
	"testing"
	"time"
)

var testSmhi *smhi
var testResponse *response

func TestGetUrl(t *testing.T) {

	testSmhi = New()
	url := testSmhi.request.getUrl("58.59", "16.18")
	expectedUrl := "http://opendata-download-metfcst.smhi.se/api/category/pmp1.5g/version/1/geopoint/lat/58.59/lon/16.18/data.json"

	if url == expectedUrl {
		return
	}
	t.Errorf("Wrong url: %s", url)
}

func TestRequest(t *testing.T) {
	smhi := New()
	testResponse = smhi.GetByLatLong("56.8769", "14.8092")
}

func TestGetMaxTempByDate(t *testing.T) {
	//testSmhi := New()
	//testResponse := testSmhi.GetByLatLong("58.59", "16.18")
	date := time.Date(2014, 8, 26, 0, 0, 0, 0, time.Local)
	temp, _ := testResponse.GetMaxTempByDate(date)
	fmt.Println(temp)
}

func TestGetMinTempByDate(t *testing.T) {
	//testSmhi := New()
	//testResponse := testSmhi.GetByLatLong("58.59", "16.18")
	date := time.Date(2014, 8, 26, 0, 0, 0, 0, time.Local)
	temp, _ := testResponse.GetMinTempByDate(date)
	fmt.Println(temp)
}
