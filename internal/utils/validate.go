package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// ValidateCSV checks that the CSV has a readable header and consistent row lengths.
func ValidateCSV(r io.Reader) error {
	// defer r.Close()

	csvReader := csv.NewReader(r)
	csvReader.TrimLeadingSpace = true
	header, err := csvReader.Read()
	if err != nil {
		return err
	} // if 
	csvReader.FieldsPerRecord = len(header)
	
	for {
		_, err := csvReader.Read()
		if err == io.EOF {
			break
		} // if 
		if err != nil {
			return err
		} // if 
	} // for 

	return nil
} // ValidateCSV

// ValidateJSON checks whether the entire file is valid JSON.
func ValidateJSON(r io.Reader) error {
	// defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	} // if 
	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON structure")
	} // if 
	return nil
} // ValidateJSON