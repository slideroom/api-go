// Package slideroomapi provides api wrappers for SlideRoom
package slideroomapi

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	apiRoot = "https://api.slideroom.com"

	defaultRequestTimeSpan = 1 * time.Minute
)

// SlideroomAPI holds information for the client (like api key and organization code)
type SlideroomAPI struct {
	baseURL             string
	apiHashKey          string
	apiAccessKey        string
	organizationCode    string
	accountEmailAddress string
	requestTimeSpan     time.Duration

	// Resources
	Export *SlideRoomResourceExport
}

// New returns an instance of a SlideroomAPI object that you can call on
func New(apiHashKey, apiAccessKey, accountEmailAddress, organizationCode string) *SlideroomAPI {
	client := &SlideroomAPI{
		baseURL:             apiRoot,
		apiHashKey:          apiHashKey,
		apiAccessKey:        apiAccessKey,
		organizationCode:    organizationCode,
		accountEmailAddress: accountEmailAddress,
		requestTimeSpan:     defaultRequestTimeSpan,
	}

	// set up the resources
	client.Export = newExportResource(client)

	return client
}

// takes a path and some params, adds in expires and sig params and returns the result
func (s *SlideroomAPI) get(path string, params url.Values) (b []byte, status int, err error) {
	res, err := s.getRaw(path, params)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	status = res.StatusCode

	return
}

func (s *SlideroomAPI) getRaw(path string, params url.Values) (res *http.Response, err error) {
	fullURL := s.generateFullURL(path, params)

	client := &http.Client{}
	request, err := http.NewRequest("GET", fullURL, nil)
	request.Header.Set("User-Agent", "SlideRoom go client ("+slideroomClientVersion+")")

	res, err = client.Do(request)

	return
}
