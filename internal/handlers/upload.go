package handlers 

import ( 
	"net/http"
	"github.com/sam8beard/csv-json-api/internal/utils"
	"fmt"
	"path/filepath"
	"encoding/json" // for constructing response 
	"path"
	"archive/zip"
	"time"
	"strings"
	"bytes"
	"io"
	"sync"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

type Response struct { 
	// ZipURL string
	SkippedFiles []string
	ConvertedFiles []string
	SkippedCounter int
	ConvertedCounter int
} // Response 

type ConvertedFile struct { 
	FileName string
	Contents []byte
} // ConvertedFile

func UploadHandler(w http.ResponseWriter, r *http.Request) { 
	response := Response{}
	response.SkippedFiles = []string{}
	response.ConvertedFiles = []string{} 
	response.SkippedCounter = 0
	response.ConvertedCounter = 0

	

	if r.Method != "POST" { 
		http.Error(w, "Only POST allowed ", http.StatusMethodNotAllowed)
		return
	} // if 
	
	// create zip writer to write files for download
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)


	
	// populate Multipart Form to retrieve file and file header (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	
	if err != nil { 
		http.Error(w, "Unable to parse form ", http.StatusBadRequest)
		return
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
					
					var convertedContents []byte
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
						convertedContents, err = utils.ConvertToJSON(fileReader)
						
						// could not convert file
						if err != nil { 
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: could not convert"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file
							continue
						} // if 
						
						// create string of file name with current suffix removed
						newFileName := strings.Replace(header.Filename, ".csv", ".json", 1)

						// write that file to the zip  
						fileWriter, err := zipWriter.Create(newFileName)
						if err != nil { 
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return 
						} // if 
						_, err = fileWriter.Write(convertedContents)
						if err != nil { 
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return 
						} // if 
						response.ConvertedFiles = append(response.ConvertedFiles, newFileName)
						response.ConvertedCounter++
						fileReader.Close()
						continue 
						
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
						convertedContents, err = utils.ConvertToCSV(fileReader)

						// could not convert file
						if err != nil { 
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: could not convert"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							// process next file
							continue
						}
						// create string of file name with current suffix removed
						newFileName := strings.Replace(header.Filename, ".json", ".csv", 1)
						
						// write that file to the zip  
						fileWriter, err := zipWriter.Create(newFileName)
						if err != nil { 
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return 
						} // if 
						_, err = fileWriter.Write(convertedContents)
						if err != nil { 
							http.Error(w, "Internal server error", http.StatusInternalServerError)
							return 
						} // if 
						response.ConvertedFiles = append(response.ConvertedFiles, newFileName)
						response.ConvertedCounter++
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
			http.Error(w, "Missing form field 'files'", http.StatusBadRequest)
		} // if
	} // if 
	
	// If non-file form fields are supplied, iterate through map and process
	// (iterate through URLs)
	if len(r.MultipartForm.Value) != 0 { 
		if r.MultipartForm.Value["urls"] != nil { 

			// declare wait group for counting tasks 
			var wg sync.WaitGroup 

			// we are expecting the amount of urls passed by the user 
			// to finish processing before continuing
			wg.Add(len(r.MultipartForm.Value["urls"]))

			// make channel to store converted files 
			convertedChan := make(chan ConvertedFile)

			for _, rawURL := range r.MultipartForm.Value["urls"] { 

				// make converted file object
				convertedFile := ConvetedFile{}

				go func() { 

					// all download, validation, and conversion logic will go here
					// fill out convertedFile object, pass to channel 
					
					wg.Done() 
				}() // routine 
				
				// attempt to download file, if not, return custom error message 
				fileReader, err := utils.DownloadFile(rawURL)
				if err != nil { 
					response.SkippedCounter++ 
					response.SkippedFiles = append(response.SkippedFiles, err.Error())
					continue
				} // if 

				fileBuf, _ := io.ReadAll(fileReader)
				
				// create readers for validation and conversion
				validationReader := io.NopCloser(bytes.NewReader(fileBuf))
				conversionReader := io.NopCloser(bytes.NewReader(fileBuf))
				
				// retrieve file type again based on revalidation 
				csvErr := utils.ValidateCSV(validationReader)
				var fileType string
				if csvErr != nil { 
					fileType = ".json"
				} else { 
					fileType = ".csv"
				} // if 
				var outputExt string
				var convertedContents []byte
				// if file type is json, convert to csv, and vice versa 
				if fileType == ".json" { 
					convertedContents, err = utils.ConvertToCSV(conversionReader)
					outputExt = ".csv"
					if err != nil { 
						msg := "URL " + rawURL + " skipped: could not convert"
						response.SkippedCounter++
						response.SkippedFiles = append(response.SkippedFiles, msg)
						continue
					} // if 
				} else if fileType == ".csv" { 
					convertedContents, err = utils.ConvertToJSON(conversionReader)
					outputExt = ".json"
					if err != nil { 
						msg := "URL " + rawURL + " skipped: could not convert"
						response.SkippedCounter++
						response.SkippedFiles = append(response.SkippedFiles, msg)
						continue
					} // if 
				} // if 
				
				// make file (determine file name here) to write converted contents to
						// maybe use raw url to construct name for file
				URLBase := path.Base(rawURL)
				t := time.Now()
				formattedTime := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d",
       				t.Year(), t.Month(), t.Day(),
       				t.Hour(), t.Minute(), t.Second())
				newFileName := URLBase + "-" + formattedTime + outputExt
				
				// write that file to the zip  
				fileWriter, err := zipWriter.Create(newFileName)
				if err != nil { 
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				} // if 
				_, err = fileWriter.Write(convertedContents)
				if err != nil { 
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				} // if 
				response.ConvertedFiles = append(response.ConvertedFiles, newFileName)
				response.ConvertedCounter++
			} // for	

			// wait for all routines to finish 
			wg.Wait()

			// close channel here
			
		} else { 
			http.Error(w, "Missing non-form field 'urls'", http.StatusBadRequest)
		} // if 
	} // if 
	
	// if converted counter is 0 respond accordingly
	if response.ConvertedCounter == 0 { 
		http.Error(w, "No files could be converted", http.StatusBadRequest)
    	return
	} // if 
	
	
	encodedResponse, err := json.MarshalIndent(response, "", "	")
	if err != nil { 
		http.Error(w, "Failed to encode log", http.StatusBadRequest)
		return
	} // if 
	// add json log/report of response 
	fileWriter, err := zipWriter.Create("log.json")
	if err != nil { 
		http.Error(w, "Failed to create log file", http.StatusBadRequest)
		return
	} // if
	_, err = fileWriter.Write(encodedResponse)
	if err != nil { 
		http.Error(w, "Failed to write log to archive", http.StatusBadRequest)
		return
	} // if 
	// close zip writer
	
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"archive.zip\"")
	w.WriteHeader(http.StatusOK)

	err = zipWriter.Close()
	if err != nil {
		http.Error(w, "Failed to finalize ZIP archive", http.StatusInternalServerError)
		return
	} // if 

	_, err = w.Write(buf.Bytes())
} // UploadHandler 