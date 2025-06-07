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
					fileReader, err := header.Open()
					if err != nil {
						fmt.Println("Error: cannot open file")
						response.SkippedCounter++
						msg := "file " + header.Filename + " skipped: cannot open file"
						response.SkippedFiles = append(response.SkippedFiles, msg)
						continue
					} // if 
					
					if fileExtension == ".csv" {
						err := utils.ValidateCSV(fileReader)
						fileReader.Close()
						if err != nil {
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: incorrect formatting"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						// Execute convert here
						fileReader, err := header.Open()
						if err != nil {
							fmt.Println("Error: cannot open file")
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: cannot open file"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						convertedFile := utils.ConvertToJSON(fileReader)
						fileReader.Close()
						continue
						/* CLOSE FILE AFTER PROCESSING */
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
							fmt.Println("Error: cannot open file")
							response.SkippedCounter++
							msg := "file " + header.Filename + " skipped: cannot open file"
							response.SkippedFiles = append(response.SkippedFiles, msg)
							continue
						} // if 
						convertedFile := utils.ConvertToCSV(fileReader)
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
				fmt.Fprintf(w, "Size of %s: %v bytes\n", header.Filename, header.Size)
			} // for 
		} else { 
			fmt.Fprintln(w, "For files, please use the field name, 'files'.")
		}
	} // if 
	
	// If non-file form fields are supplied, iterate through map and process
	if len(r.MultipartForm.Value) != 0 { 
		if r.MultipartForm.Value["urls"] != nil { 
			for _, rawUrl := range r.MultipartForm.Value["urls"] { 
				parsedUrl, err := url.Parse(rawUrl)
				if err != nil { 
					/* 
					ERROR: Add something to error report, invalid url 
					*/
					response.SkippedCounter++
					msg := "URL " + rawUrl + " skipped: could not parse"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					continue
				} // if 
				fileExtension := filepath.Ext(parsedUrl.Path)
				if fileExtension == ".csv" || fileExtension == ".json" { 
					/*
					VALID: Download then convert
					*/
				} else { 
					response.SkippedCounter++
					msg := "URL" + parsedUrl + " skipped: unsupported file type"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					continue
				} // if 
				fmt.Fprintln(w, url)
			} // for	
		} else { 
			fmt.Fprintln(w, "For urls, please use the field name, 'urls'.")
		} // if 
	} // if 
	
	/*
		Eventually, this will also be populated with the remaining response members
	*/
	
	// Encode response and write it to response writer 
	encodedResponse, err := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
} // UploadHandler 