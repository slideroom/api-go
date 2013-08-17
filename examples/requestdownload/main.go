package main

import (
	"fmt"
	"github.com/slideroom/sdk-go/sdk"
	"io/ioutil"
	"time"
)

const (
  apiKey = "Your API Key"
  loginEmail = "your@email.com"
  organizationCode = "sampleorgcode"
)

func main() {
	s := sdk.New(apiKey, loginEmail, organizationCode)

	// make a request
	//requestRes, err := s.RequestExportWithSearch("Sample", sdk.Csv, "Dallas")
	requestRes, err := s.RequestExport("Sample", sdk.Txt)

	if err != nil {
		panic(err)
	}

	// check every second until it is done
	c := time.Tick(1 * time.Second)

	for {
		downloadRes, err := s.DownloadExport(requestRes.Token)

		if err != nil {
			continue
		}

		if downloadRes.Pending == false {
			b, err := ioutil.ReadAll(*downloadRes.Export)
			if err != nil {
				panic(err)
			}

			fmt.Printf(string(b))
			return
		}

		<-c
	}
}
