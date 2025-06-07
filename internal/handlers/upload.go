package handlers 

import ( 
	"net/http"
	// "github.com/sam8beard/csv-json-api/internal/utils"
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
	response.SkippedFiles := []string{}
	response.ConvertedFiles := []string{} 
	response.ZipURL := ""
	response.FilesSkipped := 0
	response.FilesProcessed := 0

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
				file_extension := filepath.Ext(header.Filename)
				if file_extension == ".csv" || file_extension == "json" {
					
					// pass to validate then pass to convert? or
					// call validate in convert? 

				} else { 
					
					/* 
					ERROR: 

					Within response, include a skipped field that contains strings indicating which files 
					had the wrong extension -- deal with this later
					*/
					response.SkippedCounter++
					msg := "file " + header.Filename + " skipped: invalid file extension. Must be .csv or .json"
					response.SkippedFiles = append(response.SkippedFiles, msg)
					

				}
				fmt.Fprintf(w, "Size of %s: %v bytes\n", header.Filename, header.Size)
			} // for 
		} else { 
			fmt.Fprintln(w, "For files, please use the field name, 'files'.")
		}
	} // if 
	
	// If non-file form fields are supplied, iterate through map and process
	if len(r.MultipartForm.Value) != 0 { 
		if r.MultipartForm.Value["urls"] != nil { 
			for _, url := range r.MultipartForm.Value["urls"] { 
				parsedUrl, err := url.Parse(url)
				if err != nil { 
					/* 
					ERROR: Add something to error report, invalid url 
					*/
					response.SkippedCounter++
					msg := "URL " + url + " skipped: could not parse"
					response.SkippedFiles := append(response.SkippedFiles, msg)
				} // if 
				file_extension := filepath.Ext(parsedUrl)
				if file_extension == ".csv" || file_extension == ".json" { 
					/*
					VALID: Download then convert
					*/
				} else { 
					/* 
					ERROR: Add to report, wrong type of file 
					*/
					response.SkippedCounter++
					msg := "URL" + parsedUrl + " skipped: unsupported file type"
					response.SkippedFiles := append(response.SkippedFiles, parsedUrl)

				} // if 
				fmt.Fprintln(w, url)
			} // for	
		} else { 
			fmt.Fprintln(w, "For urls, please use the field name, 'urls'.")
		} // if 
	} // if 
} // UploadHandler 