package handlers 

import ( 
	"net/http"
	// "github.com/sam8beard/csv-json-api/internal/utils"
	"os"
	"fmt"
	"path/filepath"
	"net/url"
	"encoding/json" // for constructing response 
	"errors"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)
/* 

CTRL F TO SEE WHAT NEEDS TO BE CLEANED UP

*/
type Response struct { 
	ZipURL string
	SkippedFiles []string
	ConvertedFiles []string
	SkippedCounter int
	ConvertedCounter int
} // Response 

func UploadHandler(w http.ResponseWriter, r *http.Request) { 
	response := Response{}
	response.SkippedFiles = []string{}
	response.ConvertedFiles = []string{} 
	response.ZipURL = ""
	response.SkippedCounter = 0
	response.ConvertedCounter = 0

	if r.Method != "POST" { 
		http.Error(w, "Only POST allowed ", http.StatusMethodNotAllowed)
	} // if 
	
	// populate Multipart Form to retrieve file and file header (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	
	if err != nil { 
		http.Error(w, "Unable to parse form ", http.StatusBadRequest)
	} // if 

	// If file form fields are supplied, iterate through map and process
	if len(r.MultipartForm.File) != 0 {
		if r.MultipartForm.File["files"] != nil { 
			for _, header := range r.MultipartForm.File["files"] {
				fileExtension := filepath.Ext(header.Filename)
				// check file extension
				if fileExtension == ".csv" || fileExtension == ".json" {

					// attempt to open file, log if unable 
					fileReader, err := header.Open()
					if err != nil {
						response.SkippedCounter++
						msg := "file " + header.Filename + " skipped: cannot open file"
						response.SkippedFiles = append(response.SkippedFiles, msg)
						// process next file
						continue
					} // if 
					
					// if file is a csv file
					if fileExtension == ".csv" {

						// attempt to validate file
						err := utils.ValidateCSV(fileReader)
						fileReader.Close()

						// could not validate formatting of file
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: invalid or unsupported formatting"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file 
							continue
						} // if 
						
						// attempt to open file 
						fileReader, err := header.Open()

						// could not open file
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: cannot open file"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file
							continue
						} // if 

						// attempt to convert file
						convertedFile, err := utils.ConvertToJSON(fileReader)

						// could not convert file
						if err != nil { 
							response.SkippedCounter++
							msg := errors.New("file " + header.Filename + " skipped: could not convert")
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file
							continue
						} // if 
						_ = convertedFile // REMOVE WHEN CONVERT IS FINISHED
						fileReader.Close()
						continue // NOT SURE IF I NEED THIS? 
						
					// if file is a json file
					} else {

						// attempt to validate file 
						err := utils.ValidateJSON(fileReader)
						fileReader.Close()

						// could not validate formatting of file 
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: invalid or unsupported formatting"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file 
							continue
						} // if 

						// attempt to open file 
						fileReader, err = header.Open()

						// could not open file 
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: cannot open file"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file
							continue
						} // if 

						// attempt to convert file 
						convertedFile, err := utils.ConvertToCSV(fileReader)

						// could not convert file
						if err != nil { 
							response.SkippedCounter++
							msg := errors.New("file " + header.Filename + " skipped: could not convert")
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file
							continue
						}
						_ = convertedFile // REMOVE WHEN CONVERT IS FINISHED 
						fileReader.Close()
						continue
						/* CLOSE FILE AFTER PROCESSING */
					} // if
				} else { 
					response.SkippedCounter++
					msg := "file " + header.Filename + " skipped: unsupported file type"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					continue 
				} // if 
			} // for 
		} else { 
			// THIS NEEDS TO CHANGE
			fmt.Fprintln(w, "For files, please use the field name, 'files'.")
		} // if
	} // if 
	
	// If non-file form fields are supplied, iterate through map and process
	if len(r.MultipartForm.Value) != 0 { 
		if r.MultipartForm.Value["urls"] != nil { 
			for _, rawURL := range r.MultipartForm.Value["urls"] { 

				
				// if err != nil { 
				// 	response.SkippedCounter++
				// 	msg := "URL " + rawUrl + " skipped: could not parse"
				// 	response.SkippedFiles = append(response.SkippedFiles, msg)
				// 	continue
				// } // if 

				/* 

				Need to find a way to detect what kind of file reader is returned by 
				DownloadFile(). 

				This needed to decide which convert function to call

				I think we can just call both validate functions again 
				to see which file type it is

				Each validate function should return an error if incorrect file type is supplied?

				*/

				// attemp to download file, if not, return custom error message 
				fileReader, err := DownloadFile(rawURL)
				if err != nil { 
					response.SkippedCounter++ 
					response.SkippedFiles = append(response.SkippedFiles, err.Error())
					continue
				} // if 

				// check file type for conversion, if not csv, then must be json
				err := ValidateCSV(fileReader)
				if err != nil {
					// convert json file to csv
					fileContents, err := ConvertToCSV(fileReader)
					_ = err // this error doesn't need to be dealt with

					// WRITE CONVERTED FILE TO ZIP USING ZIPWRITER

				} else { 
					// convert csv file to json 
					fileContents, err := ConvertToJSON(fileReader)
					_ = err // this error doesn't need to be dealt with

					// WRITE CONVERTED FILE TO ZIP USING ZIPWRITER

				} // if 
			} // for	
		} else { 
			// THIS NEEDS TO BE CHANGED 
			fmt.Fprintln(w, "For urls, please use the field name, 'urls'.")
		} // if 
	} // if 
	
	// Encode response and write it to response writer 
	encodedResponse, err := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
} // UploadHandler 