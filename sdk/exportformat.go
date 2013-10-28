package sdk

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
func (e ExportFormat) String() string {
	switch e {
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
