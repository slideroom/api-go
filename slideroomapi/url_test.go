package slideroomapi

import (
	"testing"
)

func TestParamsSigning(t *testing.T) {
	paramsTable := []struct {
		url      string
		expected string
	}{
		{url: "http://localhost/?test=test", expected: "h/MjoGv0pzlqRW2ETn7TvYpb99Q="},
		{url: "http://localhost.localhost/?test=test", expected: "h/MjoGv0pzlqRW2ETn7TvYpb99Q="},            // domains do not matter
		{url: "http://localhost.localhost/org/123.txt?test=test", expected: "h/MjoGv0pzlqRW2ETn7TvYpb99Q="}, // neither do paths
	}

	for _, v := range paramsTable {
		got := client.generateSignature(v.url)
		if got != v.expected {
			t.Errorf("Expected: %s, Got: %s", v.expected, got)
		}
	}
}
