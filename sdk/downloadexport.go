package sdk

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// DownloadExportResponse contains server response of a download
type DownloadExportResponse struct {
	// If the download is still generating, this will be true
	Pending bool

	// The actual data of the download, if pending is false
	Export io.ReadCloser
}

func downloadExportEndpointURL() string {
	return fmt.Sprintf("export/download")
}

// DownloadExport checks a token
func (s *SlideroomAPI) DownloadExport(token int) (downloadRes *DownloadExportResponse, err error) {
	path := downloadExportEndpointURL()

	params := url.Values{}
	params.Add("token", strconv.Itoa(token))
	res, err := s.getRaw(path, params)

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
