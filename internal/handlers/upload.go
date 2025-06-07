package handlers 

import ( 
	"net/http"
	// "github.com/sam8beard/csv-json-api/internal/utils"
	"fmt"
	"path/filepath"
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

					log? 
					print statment? 
					
					*/
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
			for _, value := range r.MultipartForm.Value["urls"] { 
				/* Download file here */ 

				/*					 */
				fmt.Fprintln(w, value)
			} // for	
		} else { 
			fmt.Fprintln(w, "For urls, please use the field name, 'urls'.")
		} // if 
	} // if 
} // UploadHandler 