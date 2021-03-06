package slideroomapi

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var goodRequestResponse = `{
      "token": 123,
      "submissions": 22,
      "message": ""
    }`

func mockRequestExportArgs() (urlToHandle, exportName string, format ExportFormat) {
	exportName = "test"
	format = Csv

	urlToHandle = "/orgcode/export/request"

	return
}

func TestSuccessRequest(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockRequestExportArgs()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.URL.Query()["since"]
		if ok == true {
			t.Error("parameter 'since' should not be here")
		}

		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, goodRequestResponse)
	})

	_, err := client.Export.Request(exportName, format)

	if err != nil {
		t.Error("Expected no error in RequestExport")
	}
}

func TestSuccessRequestSince(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockRequestExportArgs()

	date := time.Now()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		since, ok := r.URL.Query()["since"]
		if ok == false || since[0] != strconv.FormatInt(date.Unix(), 10) {
			http.Error(w, "error", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, goodRequestResponse)
	})

	_, err := client.Export.RequestSince(exportName, format, date)

	if err != nil {
		t.Error("Expected no error in RequestExport")
	}
}

func TestSuccessRequestWithSearch(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockRequestExportArgs()
	searchTestName := "test search"

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		searchName, ok := r.URL.Query()["ss"]
		if ok == false || searchName[0] != searchTestName {
			http.Error(w, "error", http.StatusBadRequest)
			return
		}

		_, ok = r.URL.Query()["since"]
		if ok == true {
			t.Error("parameter 'since' should not be here")
		}

		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, goodRequestResponse)
	})

	_, err := client.Export.RequestWithSearch(exportName, format, searchTestName)

	if err != nil {
		t.Error("Expected no error in RequestExport")
	}
}

func TestSuccessRequestWithSearchSince(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockRequestExportArgs()
	searchTestName := "test search"
	date := time.Now()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		searchName, ok := r.URL.Query()["ss"]
		if ok == false || searchName[0] != searchTestName {
			http.Error(w, "error", http.StatusBadRequest)
			return
		}

		since, ok := r.URL.Query()["since"]
		if ok == false || since[0] != strconv.FormatInt(date.Unix(), 10) {
			http.Error(w, "error", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, goodRequestResponse)
	})

	_, err := client.Export.RequestWithSearchSince(exportName, format, searchTestName, date)

	if err != nil {
		t.Error("Expected no error in RequestExport")
	}
}

func TestBadRequest(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, exportName, format := mockRequestExportArgs()

	mux.HandleFunc(urlToHandle, mockErrorResponse("bad"))

	_, err := client.Export.Request(exportName, format)

	if err == nil {
		t.Error("Expected error in RequestExport")
	}
}

func mockExportDownloadArgs() (urlToHandle string, token int) {
	token = 123
	urlToHandle = "/orgcode/export/download"

	return
}

func TestSuccessDownload(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, token := mockExportDownloadArgs()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	res, err := client.Export.Download(token)

	if err != nil {
		t.Error("Expected no error in DownloadExport")
	}

	if res.Pending == true {
		t.Error("Expected pending = false")
	}
}

func TestPendingDownload(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, token := mockExportDownloadArgs()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	res, err := client.Export.Download(token)

	if err != nil {
		t.Error("Expected no error in DownloadExport")
	}

	if res.Pending == false {
		t.Error("Expected pending = true")
	}
}

func TestBadDownload(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, token := mockExportDownloadArgs()

	mux.HandleFunc(urlToHandle, mockErrorResponse("bad"))

	_, err := client.Export.Download(token)

	if err == nil {
		t.Error("Expected error in RequestExport")
	}
}
