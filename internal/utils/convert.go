package utils

import ( 
	"fmt"
	"encoding/csv"
	"io"
	"encoding/json"
	"bytes"
	"strings"
)

func ConvertToJSON(r io.ReadCloser) ([]byte, error) { 
	csvReader := csv.NewReader(r)

	// might have to remove this
	csvReader.TrimLeadingSpace = true
	var funcErr error
	// extract header for json object keys
	keys, err := csvReader.Read(); if err != nil {funcErr = err}

	// might have to remove this 
	for i, key := range keys {
        keys[i] = strings.Trim(key, "\"")
    }
	data := make([]map[string]string, 0)
	
	for { 
		row, err := csvReader.Read()
		// check if we've reached the end of the file 
		if err == io.EOF { break } else if err != nil { fmt.Println(err) }
		entry := make(map[string]string, 0)
		// for each field name, create an entry populated with its respective field value 
		for i:=0; i < len(keys); i++ { 
			entry[keys[i]] = row[i]
		} // for 
		data = append(data, entry)
	} // for 
	// Testing 
	encodedData, err := json.Marshal(data); if err != nil { funcErr = err } 
	// fmt.Println(string(encodedData))
	return encodedData, funcErr
} // ConvertToJSON

func ConvertToCSV(r io.ReadCloser) ([]byte, error) { 
	var funcErr error
	jsonDecoder := json.NewDecoder(r)
	var data []map[string]interface{}
	
	// extract json data into map slice 
	err := jsonDecoder.Decode(&data); if err != nil { funcErr = err; return nil, funcErr }
	header := make([]string, 0)
	rows := make([][]string, 0)
	
	// extract
	seen := make(map[string]bool)
	for _, obj := range data { 
		for key := range obj { 
			if !seen[key] { 
				seen[key] = true
				header = append(header, key)
			} // if 
		} // for 	
	} // for 

	for _, obj := range data { 
		row := []string{}

		for _, key := range header { 
			value := obj[key]
			// fmt.Printf("%T\n", value)
			str, ok := value.(string)
			// if value is a string, add
			if ok { 
				row = append(row, str)
			// if value doesnt exist, add empty string
			} else if value == nil {
				nilPlaceholder := ""
				row = append(row, nilPlaceholder)
			// if other type, convert to string and add
			} else { 
				stringVal := fmt.Sprint(value)
				// fmt.Printf("Value Type: %T\n", stringVal)
				row = append(row, stringVal)
			} // for 
		} // for 
		rows = append(rows, row)
	} // for 

	
	// need to use buffer as the writer passed to NewWriter because 
	// we are not writing directly to a file 
	var buffer bytes.Buffer
	csvWriter :=  csv.NewWriter(&buffer)

	// write header using Write() 
	err = csvWriter.Write(header); if err != nil {funcErr = err}
	// write rest of rows using WriteAll()
	err = csvWriter.WriteAll(rows); if err != nil {funcErr = err}
	
	csvWriter.Flush()

	byteArray := buffer.Bytes()
	// fmt.Println(string(byteArray))
	return byteArray, funcErr
} // ConvertToCSV

/* 
CSV to JSON
	•	Use csv.NewReader to read headers and rows.
	•	Treat first row as keys, map each subsequent row as a JSON object.
	•	Store objects in a slice of map[string]string or map[string]interface{}.
	•	Encode the slice with json.Marshal.

Pitfalls to handle:
	•	Missing headers
	•	Uneven field counts (already handled by validation)

⸻

JSON to CSV
	•	Use json.Decoder to decode the input as a slice of objects.
	•	Dynamically extract all keys (from first object or union of all keys).
	•	Write headers first, then loop through each object, writing row values.

Pitfalls to handle:
	•	Inconsistent key sets across objects
	•	Non-object types (like top-level arrays of strings/numbers) — reject early

⸻

Plan for Testing
	•	Create small test files for edge cases:
	•	CSV: no headers, inconsistent fields, quoted fields
	•	JSON: arrays of objects, nested values (which should be flattened or rejected)

⸻

File Cleanup and Reusability

Design both convert functions to:
	•	Accept an io.Reader (source)
	•	Possibly return an io.Reader (e.g. via bytes.Buffer)
	•	Or optionally accept an io.Writer as a target

Avoid writing directly to disk unless your handler calls for it.

⸻

Optional Enhancements Later
	•	Add support for streaming (large file conversion without full memory load)
	•	Add flags for customizing delimiter, pretty-printing JSON, etc.
*/