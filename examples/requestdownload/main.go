package main

import (
	"fmt"
	"github.com/slideroom/sdk-go/sdk"
	"io/ioutil"
	"time"
)

const (
	hashKey          = "Your Hash Key"
	accessKey        = "Your Access Key"
	loginEmail       = "sample@sample.com"
	organizationCode = "sample"
)

func main() {
	s := sdk.New(hashKey, accessKey, loginEmail, organizationCode)

	// make a request
	//requestRes, err := s.Export.RequestWithSearch("Sample", sdk.Csv, "Dallas")
	requestRes, err := s.Export.Request("Sample", sdk.Csv)

	fmt.Println(requestRes.Submissions)

	if err != nil {
		panic(err)
	}

	// check every 10 seconds until it is done
	c := time.Tick(10 * time.Second)

	for {
		<-c

		downloadRes, err := s.Export.Download(requestRes.Token)

		if err != nil {
			panic(err)
		}

		if downloadRes.Pending == false {
			b, err := ioutil.ReadAll(downloadRes.Export)
			if err != nil {
				panic(err)
			}

			fmt.Printf(string(b))
			downloadRes.Export.Close()
			return
		}
	}
}
