# A SlideRoom API Client written in Go

## Example
```go
package main

import (
	"fmt"
	"github.com/slideroom/api-go/slideroomapi"
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
	s := slideroomapi.New(hashKey, accessKey, loginEmail, organizationCode)

	// make a request
	//requestRes, err := s.Export.RequestWithSearch("Sample", slideroomapi.Csv, "Dallas")
	requestRes, err := s.Export.Request("Sample", slideroomapi.Csv)

	if err != nil {
		panic(err)
	}

	fmt.Println(requestRes.Submissions)

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
```

## Install

```bash
go install github.com/slideroom/api-go/slideroomapi
```

## Documentation

[http://godoc.org/github.com/slideroom/api-go/slideroomapi](http://godoc.org/github.com/slideroom/api-go/slideroomapi)
