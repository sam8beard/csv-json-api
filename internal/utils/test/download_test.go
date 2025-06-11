package test 

import ( 
	"fmt"
	"testing"
	"net/url"
)

func TestDownloadFile(t *testing.T) { 
	rawTestURL := "https://raw.githubusercontent.com/fivethirtyeight/data/master/nba-elo/nbaallelo.csv"
	// rawTestURL := "https://jsonplaceholder.typicode.com/posts"
	parsedUrl, err := url.Parse(rawTestURL)
	_ = err
	url_string := parsedUrl.String()
	fileReader, err := DownloadFile(url_string)
	_ = fileReader
	fmt.Println("TestDownloadFile")
} // TestDownloadFile