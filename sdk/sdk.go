// Package sdk provides api wrappers for SlideRoom
package sdk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	apiRoot = "https://review.slideroom.com/api"

	defaultRequestTimeSpan = 1 * time.Minute
)

// SlideroomAPI holds information for the client (like api key and organization code)
type SlideroomAPI struct {
	apiHashKey          string
	organizationCode    string
	accountEmailAddress string
	requestTimeSpan     time.Duration
}

// New returns an instance of a SlideroomAPI object that you can call on
func New(apiHashKey, accountEmailAddress, organizationCode string) *SlideroomAPI {
	return &SlideroomAPI{
		apiHashKey,
		organizationCode,
		accountEmailAddress,
		defaultRequestTimeSpan,
	}
}

func (this *SlideroomAPI) generateSignature(url string) string {
	lowerURL := strings.ToLower(url)
	mac := hmac.New(sha1.New, []byte(this.apiHashKey))
	mac.Write([]byte(lowerURL))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func generateURL(path string, params url.Values) string {
	return fmt.Sprintf("%s/%s?%s", apiRoot, path, params.Encode())
}

func (this *SlideroomAPI) generateFullURL(path string, params url.Values) string {
	params.Add("expires", strconv.FormatInt((time.Now().Add(this.requestTimeSpan).Unix()), 10))
	params.Add("email", this.accountEmailAddress)

	fullURL := generateURL(path, params)

	sigParams := url.Values{}
	sigParams.Add("signature", this.generateSignature(fullURL))

	fullURL = fmt.Sprintf("%s&%s", fullURL, sigParams.Encode())

	return fullURL
}

// takes a path and some params, adds in expires and sig params and returns the result
func (this *SlideroomAPI) get(path string, params url.Values) (b []byte, status int, err error) {
	res, err := this.getRaw(path, params)
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

func (this *SlideroomAPI) getRaw(path string, params url.Values) (res *http.Response, err error) {
	fullURL := this.generateFullURL(path, params)

	res, err = http.Get(fullURL)
	return
}
