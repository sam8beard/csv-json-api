package test

import ( 
	"fmt"
	"encoding/csv"
	"io"
	// "golang.org/x/exp/maps"
	// "golang.org/x/exp/slices"
	// "errors"
	"encoding/json"
	// "os"
)

func ConvertToJSON(r io.Reader) ([]byte, error) { 
	csvReader := csv.NewReader(r)
	var funcErr error
	// extract header for json object keys
	keys, err := csvReader.Read(); if err != nil {fmt.Println(err)}
	fmt.Println(keys)
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
	fmt.Println("Final data: ", data)

	encodedData, err := json.Marshal(data); if err != nil { funcErr = err } 
	fmt.Println(string(encodedData))
	return encodedData, funcErr
	/*
	CSV to JSON
	•	Use csv.NewReader to read headers and rows.
	•	Treat first row as keys, map each subsequent row as a JSON object.
	•	Store objects in a slice of map[string]string or map[string]interface{}.
	•	Encode the slice with json.Marshal.
	
	*/
} // ConvertToJSON

func ConvertToCSV(r io.Reader) ([]byte, error) { 
	return []byte{}, nil
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