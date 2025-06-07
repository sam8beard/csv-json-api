package handlers 

import ( 
	"net/http"
	// "github.com/sam8beard/csv-json-api/internal/utils"
	"fmt"
	"path/filepath"
	"net/url"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)


func UploadHandler(w http.ResponseWriter, r *http.Request) { 
	if r.Method != "POST" { 
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	} // if 

	// populate Multipart Form to retrieve file and file header 
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
					// instantiate reader and pass to convert 
				} else { 
					/* 
					NEED TO HANDLE CASE WHERE USERS UPLOAD FILE THAT IS NOT CSV OR JSON

					In response, include a skipped field that contains strings indicating which files 
					had the wrong extension -- deal with this later
					*/
					msg := header.Filename + " has invalid extension. Must be .csv or .json."

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
					Add something to error report, invalid url 
					*/
				}
				file_extension := filepath.Ext(parsedUrl)

				if file_extension == ".csv" || file_extension == ".json" { 
					/*
					Download then convert
					*/
				} else { 
					/* 
					Add to report, wrong type of file 
					*/
				}
				fmt.Fprintln(w, url)
			} // for	
		} else { 
			fmt.Fprintln(w, "For urls, please use the field name, 'urls'.")
		} // if 
	} // if 
} // UploadHandler 