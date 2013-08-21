package sdk

import (
	"fmt"
	"net/http"
	"testing"
)

var goodRequestResponse = `{
      "token": 123,
      "submissions": 22,
      "message": ""
    }`

func mockExportArgs() (urlToHandle, exportName string, format ExportFormat) {
	exportName = "test"
	format = Csv
	urlToHandle = "/" + client.requestExportEndpointURL(exportName) + "." + format.String()

	return
}

func TestSuccessRequest(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockExportArgs()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, goodRequestResponse)
	})

	_, err := client.RequestExport(exportName, format)

	if err != nil {
		t.Error("Expected no error in RequestExport")
	}
}

func TestSuccessRequestWithSearch(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockExportArgs()
	searchTestName := "test search"

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		searchName, ok := r.URL.Query()["ss"]
		if ok == false || searchName[0] != searchTestName {
			http.Error(w, "error", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, goodRequestResponse)
	})

	_, err := client.RequestExportWithSearch(exportName, format, searchTestName)

	if err != nil {
		t.Error("Expected no error in RequestExport")
	}
}

func TestBadRequest(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockExportArgs()

	mux.HandleFunc(urlToHandle, mockErrorResponse("bad"))

	_, err := client.RequestExport(exportName, format)

	if err == nil {
		t.Error("Expected error in RequestExport")
	}
}
