package slideroomapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

func (e *SlideRoomResourceExport) request(exportTitle string, format ExportFormat, savedSearchName string, since int64) (res *RequestExportResponse, err error) {
	params := url.Values{}

	if len(savedSearchName) > 0 {
		params.Add("ss", savedSearchName)
	}

	if since != 0 {
		params.Add("since", strconv.FormatInt(since, 10))
	}

	params.Add("export", exportTitle)
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

// RequestWithSearch will request an export with a format filtering the results with a saved search
//
// Parameters:
//   exportTitle - The title of a Custom Export (You can find your exports in Settings->Custom Exports)
//   format - The format of the export (see ExportFormat)
//   savedSearchName - The name of a search you have saved from review.slideroom.com
func (e *SlideRoomResourceExport) RequestWithSearch(exportTitle string, format ExportFormat, savedSearchName string) (res *RequestExportResponse, err error) {
	return e.request(exportTitle, format, savedSearchName, 0)
}

// RequestWithSearchSince will request an export with a format filtering the results with a saved search
//
// Parameters:
//   exportTitle - The title of a Custom Export (You can find your exports in Settings->Custom Exports)
//   format - The format of the export (see ExportFormat)
//   savedSearchName - The name of a search you have saved from review.slideroom.com
//   since - A Time object to only pull applications since a specific time.
func (e *SlideRoomResourceExport) RequestWithSearchSince(exportTitle string, format ExportFormat, savedSearchName string, since time.Time) (res *RequestExportResponse, err error) {
	return e.request(exportTitle, format, savedSearchName, since.Unix())
}

// RequestSince will request an export with a format filtering the results with a since
//
// Parameters:
//   exportTitle - The title of a Custom Export (You can find your exports in Settings->Custom Exports)
//   format - The format of the export (see ExportFormat)
//   since - A Time object to only pull applications since a specific time.
func (e *SlideRoomResourceExport) RequestSince(exportTitle string, format ExportFormat, since time.Time) (res *RequestExportResponse, err error) {
	return e.request(exportTitle, format, "", since.Unix())
}

// Request will request an export using the complete set of submissions
//
// Parameters:
//   exportTitle - The title of a Custom Export (You can find your exports in Settings->Custom Exports)
//   format - The format of the export (see ExportFormat)
func (e *SlideRoomResourceExport) Request(exportTitle string, format ExportFormat) (res *RequestExportResponse, err error) {
	return e.request(exportTitle, format, "", 0)
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
