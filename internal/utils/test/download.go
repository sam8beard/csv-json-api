package test

import ( 
	"fmt"
	// "os"
	"net/http"
	"errors"
	"strings"
	"net/url"
	"io"
	"path/filepath"
	"bytes"
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

	fileExtension := filepath.Ext(parsedURLString)

	// if url contains an extension of some sort and that extension is not .csv or .json,
	// return an error 
	if strings.Contains(fileExtension, ".") {
		fmt.Println("This should be printing")
		if fileExtension != ".csv" && fileExtension != ".json" { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: invalid URL type")
			return nil, funcErr
		} // if 
	} // if 

	

	// check file extension, if not .csv or .json, continue execution and delay 
	// judgement until content inspection 

	// api endpoints supplied should be supported if they return .csv or .json content
	if fileExtension != ".csv" && fileExtension != ".json" { 

		// funcErr = errors.New("URL " + parsedURLString + " skipped: invalid URL type")
		// return nil, funcErr

		// attempt to retrieve file, if error, log
		response, err := http.Get(parsedURLString)
		if err != nil || response.StatusCode != 200  { 
			// may need to do more error checking here based on empty/misleading headers 
			funcErr = errors.New("URL " + parsedURLString + " skipped: file does not exist at specified location")
		 	return nil, funcErr
		} // if 
		
		// consume stream and load content 
		data, err := io.ReadAll(response.Body)
		if err != nill { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: issue reading content of response body")
			return nil, funcErr
		} // if 
		response.Body.Close()

		// create reader for validation 
		validationReader, err := bytes.NewReader(data)
		
		// reader that will be returned 
		finalReader, err := bytes.NewReader(data)

		
	} else { 
		
		/* 

		Handle Transparent Compression Here

			// Read first few bytes of response body (512 bytes?)
			// Look for gzip indicator - 0x1F 0x88

			// If gzip wrap response body in gzip.NewReader() 
			// Else use response body as is 


		*/



		// check content type 
		contentType := response.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") { 
		
			// open reader 
			jsonReader := response.Body
			fmt.Printf("%T\n", jsonReader)
			// pass reader to validate function
			err := ValidateJSON(jsonReader)
			if err != nil { 
				funcErr = errors.New("URL " + parsedURLString + " skipped: invalid or unsupported formatting")
				return nil, funcErr
			} // if 
			jsonReader.Close()
		} else if (strings.Contains(contentType, "text/plain")) || 
		(strings.Contains(contentType, "text/csv")) || 
		(strings.Contains(contentType, "application/csv")) { 

			// open reader 
			csvReader := response.Body
			
			// pass reader to validate function
			err := ValidateCSV(csvReader)
			if err != nil { 
				funcErr = errors.New("URL " + parsedURLString  + " skipped: invalid or unsupported formatting")
				return nil, funcErr
			} // if 

		} else { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: unsupported file type")
			return nil, funcErr
		} // if 

	} // if 
	fmt.Println("Testing: ", parsedURLString)
	fmt.Println(response.Header.Get("Content-Length"))
	fmt.Println(response.Header.Get("Content-Type"))
	
	return response.Body, funcErr
} // DownloadFile