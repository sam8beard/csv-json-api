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
	
	// parsed URL String for file retrieval 
	parsedURLString := parsedURL.String()

	// extension for error checking
	fileExtension := filepath.Ext(parsedURL.Path)

	// if url contains an extension of some sort and that extension is not .csv or .json,
	// return an error 
	if fileExtension != "" && fileExtension != ".csv" && fileExtension != ".json" { 
		funcErr = errors.New("URL " + parsedURLString + " skipped: invalid URL type")
		return nil, funcErr
	} // if 
	

	// check file extension, if empty (api endpoint), .csv, or .json, continue execution and delay 
	// judgement until content inspection. 
	// api endpoints supplied should be supported if they return .csv or .json content

	// THIS CONDITION IS NOT NEEDED - REMOVE!!!!
	// We can remove because: by the time the execution gets here, we know the 
	// url is either an api endpoint, a json file, or a csv file 
	if fileExtension == "" || fileExtension == ".json" || fileExtension == ".csv" { 

		// attempt to retrieve file, if error, log
		response, err := http.Get(parsedURLString)
		if err != nil || response.StatusCode != 200  { 
			// may need to do more error checking here based on empty/misleading headers !!!
			funcErr = errors.New("URL " + parsedURLString + " skipped: issue retrieving file - " + response.Status)
		 	return nil, funcErr
		} // if 
		
		// consume stream and load content for use in testing readers 
		data, err := io.ReadAll(response.Body)
		if err != nill { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: issue reading content of response body")
			return nil, funcErr
		} // if 
		response.Body.Close()
		
		// // reader that will be returned 
		// finalReader, err := bytes.NewReader(data)
		// var finalReader io.ReadCloser
		// var validationReader io.ReadCloser
		// var gzipReader gzip.Reader

		// // reader for format validation 
		// validationReader, err := bytes.NewReader(data)
		// var finalReader io.ReadCloser
		// var validateCSVReader io.ReadCloser
		// var validateJSONReader io.ReadCloser

		// // if initialized we know the file was compressed, a gzip reader will be returned
		// var gzipReader gzip.Reader

		// if response body contents is compressed - decompress, validate, and return
		if isGzip(data) { 
			tempReader := bytes.NewReader(data)
			gzipReader := gzip.NewReader(tempReader)
		
			// temp reader for validationCSV reader
			tempReader1, err := bytes.NewReader(data)
			_ = err // MIGHT HAVE TO DEAL WITH THIS

			// temp reader for validationJSON reader
			tempReader2, err := bytes.NewReader(data)
			_ = err // MIGHT HAVE TO DEAL WITH THIS

			// temp reader for final reader
			tempReader3, err:= bytes.NewReader(data)
			_ = err // MIGHT HAVE TO DEAL WITH THIS

			validationReaderCSV, err := gzip.NewReader(tempReader1)
			if err != nil { 
				fmt.Println("detected as gzip file, but is not gzip file. REVIEW")
			} // if 

			validationReaderJSON, err := gzip.NewReader(tempReader2)
			_ = err // dont have to deal with, above error will trigger if not gzip

			finalReader, err := gzip.NewReader(tempReader3)
			_ = err // same as above
			

			// validate file, at least one of these will be nil
			// if both are nill, return error 
			
			csvErr := ValidateCSV(validationReaderCSV)
			jsonErr := ValidateJSON(validationReaderJSON)

			// file was of type .csv, but formatting was invalid 
			if fileExtension == ".csv" && csvErr != nil { 
				
			// file was of type json, but formatting was invalid 
			} else if fileExtension == ".json" && jsonErr != nil { 
				
			// file returned by endpoint was neither json or csv 
			} else if fileExtension == "" && jsonErr != nil && csvErr != nill { 
				
			} // if 

			
			
		// response is not compressed - validate and return 
		} else { 

		} // if 
		
		
		/* 

		Handle Transparent Compression Here

			// Read first few bytes of response body (512 bytes?)
			// Look for gzip indicator - 0x1F 0x88

			// If gzip wrap response body in gzip.NewReader() 
			// Else use response body as is 


		*/
		
		
	
	} else  { 
		
		


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

func isGzip(data []byte) bool { 
	return len(data) > 2 && data[0] == 0x1F && data[1] == 0x8B
} // isGzip