package test

import ( 
	"fmt"
	"encoding/csv"
	"io"
	"encoding/json"
	"slices"
	"bytes"
	// "os"
)

func ConvertToJSON(r io.Reader) ([]byte, error) { 
	csvReader := csv.NewReader(r)
	var funcErr error
	// extract header for json object keys
	keys, err := csvReader.Read(); if err != nil {funcErr = err}
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
	return encodedData, funcErr
} // ConvertToJSON

func ConvertToCSV(r io.Reader) ([]byte, error) { 
	var funcErr error
	jsonDecoder := json.NewDecoder(r)
	var data []map[string]interface{}
	
	// extract json data into map slice 
	err := jsonDecoder.Decode(&data); if err != nil { funcErr = err }
	header := make([]string, 0)
	rows := make([][]string, 0)
	// _ = keys
	
	// extract keys
	for _, obj := range data { 
		row := []string{}

		// extract keys
		for key := range obj { 
			if !slices.Contains(header, key) { header = append(header, key) }
		} // for 	

		// for each key in header, extract corresponding value in current
		// object and append it to a row 
		for _, key := range header { 
			value := obj[key]
			str, ok := value.(string)
			if ok { 
				row = append(row, str)
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
	fmt.Println(string(byteArray))
	
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