package epa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestClient_GetAirQualityForecast(t *testing.T) {
	setup()
	defer teardown()

	testdata, _ := ioutil.ReadFile(fmt.Sprintf("./testdata/%v.json", ResourceIDAirQualityForecast))

	mux.HandleFunc(fmt.Sprintf("/webapi/api/rest/datastore/%v", ResourceIDAirQualityForecast), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"sort":   "PublishTime",
			"offset": "0",
			"limit":  "10",
		})

		w.WriteHeader(http.StatusOK)
		w.Write(testdata)
	})

	options := url.Values{}
	options.Set("sort", "PublishTime")
	options.Set("offset", "0")
	options.Set("limit", "10")

	got, _, err := client.GetAirQualityForecast(context.Background(), options)
	if err != nil {
		t.Errorf("GetAirQualityForecast returned error: %v", err)
	}

	restored, _ := json.Marshal(got)

	areEqual, err := areEqualJSON(testdata, restored)
	if !areEqual {
		t.Error("GetAirQualityForecast testdata and restored are not equal")
	}
}
