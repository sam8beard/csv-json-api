package handlers 

import ( 
	"net/http"
	// "github.com/sam8beard/csv-json-api/internal/utils"
	"fmt"
	"path/filepath"
	"encoding/json" // for constructing response 
	"path"
	"archive/zip"
	"time"
	"strings"
	"log"
	"bytes"
	"io"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)
/* 

CTRL F TO SEE WHAT NEEDS TO BE CLEANED UP

*/
type Response struct { 
	// ZipURL string
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
		return
	} // if 
	
	// create zip writer to write files for download
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	
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
			for _, rawURL := range r.MultipartForm.Value["urls"] { 
				// attempt to download file, if not, return custom error message 
				fileReader, err := DownloadFile(rawURL)
				if err != nil { 
					response.SkippedCounter++ 
					response.SkippedFiles = append(response.SkippedFiles, err.Error())
					continue
				} // if 

				fileBuf, _ := io.ReadAll(fileReader)
				
				// create readers for validation and conversion
				validationReader := bytes.NewReader(fileBuf)
				conversionReader := bytes.NewReader(fileBuf)
				
				// retrieve file type again based on revalidation 
				csvErr := utils.ValidateCSV(validationReader)
				var fileType string
				if csvErr != nil { 
					fileType = ".json"
				} else { 
					fileType = ".csv"
				} // if 


				// if file type is json, convert to csv, and vice versa 
				if fileType == ".json" { 
					convertedContents, err = utils.ConvertToCSV(conversionReader)
					if err != nil { 
						msg := "URL " + rawURL + " skipped: could not convert"
						response.SkippedCounter++
						response.SkippedFiles = append(response.SkippedFiles, msg)
						continue
					} // if 
				} else if fileType == ".csv" { 
					convertedContents, err = utils.ConvertToJSON(conversionReader)
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
				newFileName := URLBase + "-" + formattedTime + fileType
				
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

	_, err = w.Write(buf.Bytes())
	if err != nil { 
		log.Println("Failed to write archive to response:", err)
		return
	} // if 
} // UploadHandler 