package handlers 

import ( 
	"net/http"
	"github.com/sam8beard/csv-json-api/utils"
	// "fmt"
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
				fileReader, err := header.Open(); if err != nil { fmt.Println(err)}
				contents, err := io.ReadAll(fileReader); if err != nil { fmt.Println(err)}
				_ = fileReader
				_ = contents 
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