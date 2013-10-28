package sdk

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// SlideRoomResourceExport contains functions for handling the export resource
type SlideRoomResourceExport struct {
	client *SlideroomAPI
}

func newExportResource(client *SlideroomAPI) *SlideRoomResourceExport {
	return &SlideRoomResourceExport{
		client: client,
	}
}

// RequestExportResponse holds values about the Request/Download
type RequestExportResponse struct {
	// the token of the download that you can use to download
	Token int `json:"token"`

	// The number of submissions involved in this request
	Submissions int `json:"submissions"`

	// Any message that the server wants to send
	Message string `json:"message"`
}

// RequestWithSearch will request an export with a format filtering the results with a saved search
// You can find your export names in Settings->Custom Exports (use the title)
func (e *SlideRoomResourceExport) RequestWithSearch(exportName string, format ExportFormat, search string) (res *RequestExportResponse, err error) {
	params := url.Values{}

	if len(search) > 0 {
		params.Add("ss", search)
	}

	params.Add("export", exportName)
	params.Add("format", format.String())

	b, status, err := e.client.get("export/request", params)
	if err != nil {
		return
	}

	if status != http.StatusAccepted {
		return nil, parseErrorFromBytes(b)
	}

	err = json.Unmarshal(b, &res)
	return
}

// Request will request an export using the complete set of submissions
// You can find your export names in Settings->Custom Exports (use the title)
func (e *SlideRoomResourceExport) Request(exportName string, format ExportFormat) (res *RequestExportResponse, err error) {
	return e.RequestWithSearch(exportName, format, "")
}

// DownloadExportResponse contains server response of a download
type DownloadExportResponse struct {
	// If the download is still generating, this will be true
	Pending bool

	// The actual data of the download, if pending is false
	Export io.ReadCloser
}

// Download checks a token
func (e *SlideRoomResourceExport) Download(token int) (downloadRes *DownloadExportResponse, err error) {
	params := url.Values{}
	params.Add("token", strconv.Itoa(token))
	res, err := e.client.getRaw("export/download", params)

	if err != nil {
		return
	}

	switch res.StatusCode {
	case http.StatusAccepted:
		downloadRes = &DownloadExportResponse{
			Pending: true,
		}

	case http.StatusOK:
		downloadRes = &DownloadExportResponse{
			Pending: false,
			Export:  res.Body,
		}

	default:
		downloadRes = nil
		err = parseErrorFromResponse(res)
	}

	return
}
