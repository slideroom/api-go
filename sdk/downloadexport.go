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
	Export *io.ReadCloser
}

func (this *SlideroomAPI) downloadExportEndpointURL(token int) string {
	return fmt.Sprintf("export/%s/%s", this.organizationCode, strconv.Itoa(token))
}

// DownloadExport checks a token
func (this *SlideroomAPI) DownloadExport(token int) (downloadRes *DownloadExportResponse, err error) {
	path := this.downloadExportEndpointURL(token)
	res, err := this.getRaw(path, url.Values{})

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
			Export:  &res.Body,
		}

	default:
		downloadRes = nil
		err = parseErrorFromResponse(res)
	}

	return
}
