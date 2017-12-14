package epa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Apple Music client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a every9d.Client that is configured to talk to that test server.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// EVERY8D client configured to use test server
	client = NewClient("token", nil)
	u, _ := url.Parse(server.URL)
	client.BaseURL = u
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method is %v, want %v", got, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
	r.Form.Del("token") // Remove token
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters is %v, want %v", got, want)
	}
}

func areEqualJSON(j1, j2 []byte) (bool, error) {
	var v1 interface{}
	var v2 interface{}

	var err error
	err = json.Unmarshal(j1, &v1)
	if err != nil {
		return false, fmt.Errorf("Unmarshal JSON 1 error: %v", err)
	}
	err = json.Unmarshal(j2, &v2)
	if err != nil {
		return false, fmt.Errorf("Unmarshal JSON 2 error: %v", err)
	}

	return reflect.DeepEqual(v1, v2), nil
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusInternalServerError,
		Body:       ioutil.NopCloser(strings.NewReader(`{"Message": "An error has occurred."}`)),
	}

	err := CheckResponse(res).(*ErrorResponse)
	if err == nil {
		t.Errorf("Expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "An error has occurred.",
	}
	if !reflect.DeepEqual(err, want) {
		t.Errorf("Error = %#v, want %#v", err, want)
	}
}
