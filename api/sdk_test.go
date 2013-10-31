package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *SlideroomAPI
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = New("apikey", "secretkey", "email@email.com", "orgcode")
	url, _ := url.Parse(server.URL)
	client.baseURL = url.String()
}

func teardown() {
	server.Close()
}

func mockErrorResponse(errorMessage string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":{"message":"`+errorMessage+`"}}`)
	}
}
