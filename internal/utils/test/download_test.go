package test 

import ( 
	"fmt"
	"testing"
	"net/url"
	"io"
	// "compress/gzip"
)

func TestDownloadFile(t *testing.T) { 
	// rawTestURL := "https://raw.githubusercontent.com/fivethirtyeight/data/master/nba-elo/nbaallelo.csv"
	// rawTestURL := "https://file.com"
	// rawTestURL := "https://jsonplaceholder.typicode.com/posts"
	// rawTestURL := "https://people.sc.fsu.edu/~jburkardt/data/csv/hw_200.csv"
	rawTestURL := "http://api.open-notify.org/iss-now.json"
	parsedUrl, err := url.Parse(rawTestURL)
	_ = err
	url_string := parsedUrl.String()
	fileReader, downloadErr := DownloadFile(url_string)
	fmt.Printf("%T\n", fileReader)
	fmt.Println(downloadErr)
	fmt.Println("TestDownloadFile")
	contents, err := io.ReadAll(fileReader)
	fmt.Println(string(contents))
	// gzipReader, err := gzip.NewReader(fileReader)
	// fmt.Printf("%T\n", gzipReader)

	// fileContents, err := io.ReadAll(fileReader)
	// fmt.Printf("%T\n", fileContents)
	// if err != nil { 
	// 	// fmt.Println(err)
	// } // if 
	// fmt.Println(string(fileContents))
} // TestDownloadFile