package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func signParams(params url.Values, key string) string {
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k+"="+params.Get(k))
	}

	sort.Strings(keys)

	strToGenerateFrom := strings.Join(keys, "\n")
	strToGenerateFrom = strings.ToLower(strToGenerateFrom)

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(strToGenerateFrom))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// take a url, get the params and generate a sig based on the params
func (s *SlideroomAPI) generateSignature(endpointURL string) string {
	parsedURL, _ := url.Parse(endpointURL)
	params := parsedURL.Query()

	params.Add("access-key", s.apiAccessKey)

	return signParams(params, s.apiHashKey)
}

func (s *SlideroomAPI) normalizeURL(path string, params url.Values) string {
	normalURL := fmt.Sprintf("%s/%s", s.organizationCode, path)

	if len(params) > 0 {
		normalURL = fmt.Sprintf("%s?%s", normalURL, params.Encode())
	}

	return normalURL
}

func (s *SlideroomAPI) generateURL(path string, params url.Values) string {
	return fmt.Sprintf("%s/%s", s.baseURL, s.normalizeURL(path, params))
}

func (s *SlideroomAPI) generateFullURL(path string, params url.Values) string {
	params.Add("expires", strconv.FormatInt((time.Now().Add(s.requestTimeSpan).Unix()), 10))
	params.Add("email", s.accountEmailAddress)

	fullURL := s.generateURL(path, params)

	sigParams := url.Values{}
	sigParams.Add("signature", s.generateSignature(fullURL))

	fullURL = fmt.Sprintf("%s&%s", fullURL, sigParams.Encode())

	return fullURL
}
