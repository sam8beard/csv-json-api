package test

import ( 
	"fmt"
	// "os"
	"net/http"
	"errors"
	"strings"
	"io"
)

/* 
Handle url verification logic here

Cases to account for: 

1. Unparsable URLS 
2. Get files using http request -> verify file by... 
	1. Checking content type AND THEN 
	2. Call Validate[File Type] function 

*/ 
func DownloadFile(url string) (io.ReadCloser, error) { 
	var funcErr error 
	response, err := http.Get(url)

	// file not found at url given
	if err != nil { 
		 funcErr = errors.New("URL " + url + " skipped: file does not exist at specified location")
		 return nil, funcErr
	} // if 
	
	// parse URL and check error, return correct response 

	// check content-type, if passes -> call validate
	// if validate doesnt return error -> we are OK to convert 
	
	// Validate file formatting 
	contentType := response.Header.Get("Content-Type")
	fmt.Println(contentType)
	
	if contentType strings.Contains("application/json")

	// if (filePath.Ext(url) == ".csv") { 
	// 	fmt.Println("This is a csv file")
	// } else { 
	// 	fmt.Println("This is a json file")
	// }
	// fmt.Printf("%T\n", response)
	// fmt.Printf("%T\n", response.Body)

	// contents, err := io.ReadAll(response.Body)
	// _ = err
	// fmt.Printf("%T\n", contents)
	// // fmt.Println(string(contents))
	
	// Return response.Body (this is the Reader for the file)

	/* 
	Return error if file location does not exist 

	Already checking in upload.go if:
	1. URL is parseable 
	2. URL is .csv or .json type 

	Should return error statement in the same format as that in upload 
	Ex: "URL [file url] skipped: file does not exist at specified location"
	*/ 
	
	return response.Body, funcErr
} // DownloadFile