package test

import ( 
	"fmt"
	// "os"
	"net/http"
	"errors"
	"io"
)

func DownloadFile(url string) (io.ReadCloser, error) { 
	var funcErr error 
	response, err := http.Get(url)
	if err != nil { 
		 funcErr = errors.New("URL" + url + "skipped: file does not exist at specified location")
		 return nil, funcErr
	} // if 

	fmt.Printf("%T\n", response)
	fmt.Printf("%T\n", response.Body)

	contents, err := io.ReadAll(response.Body)
	_ = err
	fmt.Printf("%T\n", contents)


	// Return response.Body (this is the Reader for the file)
	
	/* 
	Return error if file location does not exist 

	Already checking in upload.go if:
	1. URL is parseable 
	2. URL is .csv or .json type 

	Should return error statement in the same format as that in upload 
	Ex: "URL [file url] skipped: file does not exist at specified location"
	*/ 

	return nil, funcErr
} // DownloadFile