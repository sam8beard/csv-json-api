package test 

import ( 
	"fmt"
	"testing"
	"net/url"
	"io"
	// "compress/gzip"
)

func TestDownloadFile(t *testing.T) { 
	// rawTestURL := "https://raw.githubusercontent.com/fivethirtyeight/data/master/nba-elo/nbaallelo.csv" // WORKING - SHOULD RETURN READER
	// rawTestURL := "https://file.com" // WORKING - SHOULD RETURN ERROR
	rawTestURL := "https://jsonplaceholder.typicode.com/posts" // WORKING - SHOULD RETURN READER
	// rawTestURL := "https://people.sc.fsu.edu/~jburkardt/data/csv/hw_200.csv" // WORKING - SHOULD RETURN ERROR 
	// rawTestURL := "http://api.open-notify.org/iss-now.json" // WORKING - SHOULD RETURN READER
	parsedUrl, err := url.Parse(rawTestURL)
	_ = err
	url_string := parsedUrl.String()
	fileReader, downloadErr := DownloadFile(url_string)

	if downloadErr !=  nil { 
		t.Log("Error detected\n")
		fmt.Println(downloadErr)
		return
	} // if 
	fmt.Printf("%T\n", fileReader)
	fmt.Println(downloadErr)
	fmt.Println("TestDownloadFile")
	contents, err := io.ReadAll(fileReader)
	fmt.Println(string(contents))

} // TestDownloadFile