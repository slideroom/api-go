package sdk

import (
	"net/http"
	"net/url"
	"testing"
)

func mockDownloadArgs() (urlToHandle string, token int) {
	token = 123

	urlToHandle = "/" + client.normalizeURL(downloadExportEndpointURL(), url.Values{})

	return
}

func TestSuccessDownload(t *testing.T) {
	setup()
	defer teardown()

	urlToHandle, token := mockDownloadArgs()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	res, err := client.DownloadExport(token)

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

	urlToHandle, token := mockDownloadArgs()

	mux.HandleFunc(urlToHandle, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	res, err := client.DownloadExport(token)

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

	urlToHandle, token := mockDownloadArgs()

	mux.HandleFunc(urlToHandle, mockErrorResponse("bad"))

	_, err := client.DownloadExport(token)

	if err == nil {
		t.Error("Expected error in RequestExport")
	}
}
