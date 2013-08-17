package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ExportFormat holds the formats that a download can be in
type ExportFormat int

const (
	// Csv = Comma seperated values
	Csv ExportFormat = iota

	// Tsv = Tab seperated values
	Tsv

	// Txt = text file
	Txt

	// Xlsx = excel format
	Xlsx
)

// converts an ExportFormat to an extension
func (this ExportFormat) String() string {
	switch this {
	case Csv:
		return "csv"

	case Tsv:
		return "tsv"

	case Txt:
		return "txt"

	case Xlsx:
		return "xlsx"

	default:
		return "csv"
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

func (this *SlideroomAPI) requestExportEndpointURL(resource string) string {
	return fmt.Sprintf("export/%s/%s", this.organizationCode, resource)
}

// RequestExportWithSearch will request an export with a format filtering the results with a saved search
// You can find your export names in Settings->Custom Exports (use the title)
func (this *SlideroomAPI) RequestExportWithSearch(exportName string, format ExportFormat, search string) (res *RequestExportResponse, err error) {
	params := url.Values{}

	if len(search) > 0 {
		params.Add("ss", search)
	}

	resource := fmt.Sprintf("%s.%s", exportName, format)
	b, status, err := this.get(this.requestExportEndpointURL(resource), params)
	if err != nil {
		return
	}

	if status != http.StatusAccepted {
		return nil, parseErrorFromBytes(b)
	}

	err = json.Unmarshal(b, &res)
	return
}

// RequestExport will request an export using the complete set of submissions
// You can find your export names in Settings->Custom Exports (use the title)
func (this *SlideroomAPI) RequestExport(exportName string, format ExportFormat) (res *RequestExportResponse, err error) {
	return this.RequestExportWithSearch(exportName, format, "")
}
