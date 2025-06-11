package test

import ( 
	"fmt"
	// "os"
	"net/http"
	"errors"
	// "strings"
	"net/url"
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
func DownloadFile(rawURL string) (io.ReadCloser, error) { 
	var funcErr error 
	
	// parse URL and check error
	parsedURL, err := url.Parse(rawURL)
	if err != nil { 
		funcErr = errors.New("URL " + rawURL+ " skipped: could not parse")
		return nil, funcErr
	} // if 

	// check if file exist at url
	parsedURLString := parsedURL.String()
	response, err := http.Get(parsedURLString)
	if err != nil { 
		 funcErr = errors.New("URL " + parsedURLString + " skipped: file does not exist at specified location")
		 return nil, funcErr
	} // if 

	// check content type 
	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") { 
		// validate json 
	} else if (strings.Contains(contentType, "text/plain")) || 
	(strings.Contains(contentType, "text/csv")) || 
	(strings.Contains(contentType, "application/csv")) { 
		// validate csv 
	} else { 
		funcErr = errors.New("URL" + parsedURLString + " skipped: unsupported file type")
		return nil, funcErr
	}
	
	
	
	// check content-type, if passes -> call validate
	// if validate doesnt return error -> we are OK to convert 

	// Validate file formatting 
	
	// fmt.Println(contentType)
	
	// if contentType strings.Contains("application/json")

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