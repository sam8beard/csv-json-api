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
	var finalReader io.ReadCloser
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

	// attempt to retrieve file, if error, log
	response, err := http.Get(parsedURLString)
	if err != nil || response.StatusCode != 200  { 
		// may need to do more error checking here based on empty/misleading headers !!!
		funcErr = errors.New("URL " + parsedURLString + " skipped: issue retrieving file - " + response.Status)
		return nil, funcErr
	} // if 

	// consume stream and load content for use in testing readers 
	data, err := io.ReadAll(response.Body)
	if err != nil { 
		funcErr = errors.New("URL " + parsedURLString + " skipped: issue reading content of response body")
		return nil, funcErr
	} // if 
	response.Body.Close()

	// if response body contents is compressed - decompress, validate, and return
	if isGzip(data) { 

		gzipReader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: marked as gzip but not compressed")
			return nil, funcErr
		} // if 
		
		// validate that there isnt trailing uncompressed content in flie, return an error if so
		cleanFile, err := gzip.NewReader(bytes.NewReader(data))
		_ = err
		buf := make([]byte, 512)
		_, err = cleanFile.Read(buf)
		if err != nil { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: compressed file is malformed")
			return nil, funcErr
		} // if 

		validationReaderCSV, err := gzip.NewReader(bytes.NewReader(data))
		_ = err // dont have to deal with, above error will trigger if not gzip

		validationReaderJSON, err := gzip.NewReader(bytes.NewReader(data))
		_ = err 

		
		// validate file, at least one of these will be nil
		// if both are nil, return error 
		
		csvErr := ValidateCSV(validationReaderCSV)
		jsonErr := ValidateJSON(validationReaderJSON)

		// file was of type csv, but formatting was invalid 
		if fileExtension == ".csv" && csvErr != nil { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: invalid or unsupported formatting")
			return nil, funcErr
		// file was of type json, but formatting was invalid 
		} else if fileExtension == ".json" && jsonErr != nil { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: invalid or unsupported formatting")
			return nil, funcErr
		// file returned by endpoint was neither json or csv 
		} else if fileExtension == "" && jsonErr != nil && csvErr != nil { 
			funcErr = errors.New("URL " + parsedURLString + " skipped: invalid file type")
			return nil, funcErr
		} // if 
		finalReader = gzipReader
		
	// response is not compressed - validate and return 
	} else { 
		// TO DO
		// handle non compressed file here 
			
	} // if 

	return finalReader, nil
			
} // DownloadFile

func isGzip(data []byte) bool { 
	return len(data) > 2 && data[0] == 0x1F && data[1] == 0x8B
} // isGzip