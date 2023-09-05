package framework

import (
	"net/url"
)

var (
	DevURL  *url.URL
	ProdURL *url.URL
)

func init() {
	var err error

	DevURL, err = new(url.URL).Parse("http://localhost:2040/")
	if err != nil {
		panic(err.Error())
	}

	ProdURL, err = new(url.URL).Parse("https://auth.agoradesecrivains.net/")
	if err != nil {
		panic(err.Error())
	}
}
