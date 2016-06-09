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
	expectedUrl := "http://opendata-download-metfcst.smhi.se/api/category/pmp2g/version/2/geotype/point/lon/16.18/lat/58.59/data.json"

	if url == expectedUrl {
		return
	}
	t.Errorf("Wrong url: %s", url)
}

func TestRequest(t *testing.T) {
	smhi := New()
	testResponse = smhi.GetByLatLong("56.8769", "14.8092")
	fmt.Printf("%#v\n", testResponse)
}

func TestGetMaxTempByDate(t *testing.T) {
	date := time.Now()
	temp, _ := testResponse.GetMaxTempByDate(date)
	fmt.Println(temp)
}

func TestGetMinTempByDate(t *testing.T) {
	date := time.Now()
	temp, _ := testResponse.GetMinTempByDate(date)
	fmt.Println(temp)
}
