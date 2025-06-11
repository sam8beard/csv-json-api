package handlers 

import ( 
	"net/http"
	// "github.com/sam8beard/csv-json-api/internal/utils"
	"os"
	"fmt"
	"path/filepath"
	"net/url"
	"encoding/json" // for constructing response 
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
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	} // if 
	
	// populate Multipart Form to retrieve file and file header (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	
	if err != nil { 
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
	} // if 

	// If file form fields are supplied, iterate through map and process
	if len(r.MultipartForm.File) != 0 {
		if r.MultipartForm.File["files"] != nil { 
			for _, header := range r.MultipartForm.File["files"] {
				fileExtension := filepath.Ext(header.Filename)
				if fileExtension == ".csv" || fileExtension == ".json" {

					// attempt to open file, log if unable 
					fileReader, err := header.Open()
					if err != nil {
						response.SkippedCounter++
						msg := "file " + header.Filename + " skipped: cannot open file"
						response.SkippedFiles = append(response.SkippedFiles, msg)
						continue
					} // if 
					
					// if file is a csv file
					if fileExtension == ".csv" {
						err := utils.ValidateCSV(fileReader)
						fileReader.Close()
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: invalid/unsupported formatting"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						// Execute convert here
						fileReader, err := header.Open()
						if err != nil {
							fmt.Println("Error: cannot open file") // THIS NEEDS TO CHANGE 
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: cannot open file"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						convertedFile := utils.ConvertToJSON(fileReader)
						_ = convertedFile // REMOVE WHEN CONVERT IS FINISHED
						fileReader.Close()
						continue
					} else {
						err := utils.ValidateJSON(fileReader)
						fileReader.Close()
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: incorrect formatting"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						// Execute convert here
						fileReader, err = header.Open()
						if err != nil {
							fmt.Println("Error: cannot open file") // THIS NEEDS TO CHANGE
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: cannot open file"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						convertedFile := utils.ConvertToCSV(fileReader)
						_ = convertedFile // REMOVE WHEN CONVERT IS FINISHED 
						fileReader.Close()
						continue
						/* CLOSE FILE AFTER PROCESSING */
					} // if
				} else { 
					response.SkippedCounter++
					msg := "file " + header.Filename + " skipped: invalid file extension. Must be .csv or .json"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					continue 
				} // if 
			} // for 
		} else { 
			// THIS NEEDS TO CHANGE
			fmt.Fprintln(w, "For files, please use the field name, 'files'.")
		}
	} // if 
	
	// If non-file form fields are supplied, iterate through map and process
	if len(r.MultipartForm.Value) != 0 { 
		if r.MultipartForm.Value["urls"] != nil { 
			for _, rawUrl := range r.MultipartForm.Value["urls"] { 

				// add this case to download.go
				parsedUrl, err := url.Parse(rawUrl)
				if err != nil { 
					response.SkippedCounter++
					msg := "URL " + rawUrl + " skipped: could not parse"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					continue
				} // if 
				fileExtension := filepath.Ext(parsedUrl.Path)

				// add this case to download go (will need to be modified because of csv files returning text/plain sometimes)
				// text/plain might not be csv, but thats okay -> we can call ValidateCSV on the reader returned 
				if fileExtension == ".csv" || fileExtension == ".json" { 
					/*
					VALID: Download then convert
					*/
					// Still need to validate whether or not 
					// content of files are correctly formatted 
					

					/* 
					1. Download file
					2. Validate file 
					3. Convert file 
					*/
				} else { 
					response.SkippedCounter++
					msg := "URL" + parsedUrl + " skipped: unsupported file type"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					continue
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