package slideroomapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Error is a wrapper object that holds server error reposnse
type Error struct {
	Err errorObject `json:"error"`
}

type errorObject struct {
	Message string `json:"message"`
}

// Error returns the server error message
func (this Error) Error() string {
	return this.Err.Message
}

func parseErrorFromResponse(res *http.Response) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return parseErrorFromBytes(b)
}

func parseErrorFromBytes(b []byte) (err error) {
	errorResponse := &Error{}
	err = json.Unmarshal(b, &errorResponse)
	if err == nil {
		err = errorResponse
	}

	return
}
