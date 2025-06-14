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
	"path"
	"archive/zip"
	"time"
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

	// create zip archive to return for download
	zipArchive, err := os.Create("archive.zip")
	if err != nil { 
		panic(err)
	} // if 
	defer archive.Close()

	zipWriter := zip.NewWriter(zipArchive)
	defer zipWriter.Close() 

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
				
				// retrieve file type again based on revalidation 
				csvErr := ValidateCSV(fileReader)
				var fileType string
				if csvErr != nil { 
					fileType = ".json"
				} else { 
					fileType = ".csv"
				} // if 

				var convertedContents []byte

				// if file type is json, convert to csv, and vice versa 
				if fileType == ".json" { 
					convertedContents, err := ConvertToCSV(fileReader)
					if err != nil { 
						msg := "URL " + rawURL + " skipped: could not convert"
						response.SkippedCounter++
						response.SkippedFiles = append(response.SkippedFiles, msg)
					} // if 
				} else if fileType == ".csv" { 
					convertedContents, err := ConvertToJSON(fileReader)
					if err != nil { 
						msg := "URL " + rawURL + " skipped: could not convert"
						response.SkippedCounter++
						response.SkippedFiles = append(response.SkippedFiles, msg)
					} // if 
				} // if 
				
				// make file (determine file name here) to write converted contents to
						// maybe use raw url to construct name for file
				URLBase := path.Base(rawURL)
				t := time.Now()
				formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
       				t.Year(), t.Month(), t.Day(),
       				t.Hour(), t.Minute(), t.Second())
				newFileName := URLBase + "_" + formattedTime + fileType

				// write that file to the zip 
				
				
				
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