package test 

import ( 
	"fmt"
	"testing"
	"net/url"
	// "io"
	"compress/gzip"
)

func TestDownloadFile(t *testing.T) { 
	// rawTestURL := "https://raw.githubusercontent.com/fivethirtyeight/data/master/nba-elo/nbaallelo.cs"
	// rawTestURL := "https://file.com"
	rawTestURL := "https://jsonplaceholder.typicode.com/posts"
	parsedUrl, err := url.Parse(rawTestURL)
	_ = err
	url_string := parsedUrl.String()
	fileReader, downloadErr := DownloadFile(url_string)
	fmt.Printf("%T\n", fileReader)
	fmt.Println(downloadErr)
	fmt.Println("TestDownloadFile")
	gzipReader, err := gzip.NewReader(fileReader)
	fmt.Printf("%T\n", gzipReader)

	// fileContents, err := io.ReadAll(gzipReader)
	if err != nil { 
		fmt.Println(err)
	} // if 
	// fmt.Println(fileContents)
} // TestDownloadFile