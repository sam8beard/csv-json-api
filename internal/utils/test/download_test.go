package test 

import ( 
	"fmt"
	"testing"
	"net/url"
)

func TestDownloadFile(t *testing.T) { 
	// rawTestURL := "https://raw.githubusercontent.com/fivethirtyeight/data/master/nba-elo/nbaallelo.cs"
	// rawTestURL := "https://file.com"
	rawTestURL := "https://jsonplaceholder.typicode.com/posts"
	parsedUrl, err := url.Parse(rawTestURL)
	_ = err
	url_string := parsedUrl.String()
	fileReader, downloadErr := DownloadFile(url_string)
	fmt.Println(downloadErr)
	fmt.Println("TestDownloadFile")
	_ = fileReader
	
} // TestDownloadFile