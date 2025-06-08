package test

import ( 
	"io"
	"encoding/json"
	"encoding/csv"
	"io/ioutil"
	"fmt"
)

/* 
Conditions to check for csv files: 
1.	File is parseable (no malformed CSV structure)
2.	Header exists (optional depending on your rules)
3.	All rows have same number of columns
4.	No empty rows / fields
5.	Required columns exist (e.g., “id”, “name”, etc.)
6.	Field value formats (e.g., emails, numbers)
*/
func ValidateCSV(r io.Reader) error { 
	// fmt.Println("Testing")
	csvReader := csv.NewReader(r) 
	// Read file header first
	header, err := csvReader.Read(); if err != nil { return err }
	csvReader.FieldsPerRecord = len(header)
	// fmt.Println("File header: ", header)
	for { 
		_, err := csvReader.Read()
		if err == io.EOF { break } else if err != nil { return err }
		// For testing
		// if err == nil { fmt.Println(row)}
	}
	return err
} // ValidateCSV

/*
I know that encoding/json literally has a built in function to 
validate json objects and this is redundant in every way, 
but I'm instantiating my own function for my sanity.

I'm just wrapping the call to json.Validate in this function 
and returning an error instead of bool.
*/
func ValidateJSON(r io.Reader) error {
	fileContents, readErr := ioutil.ReadAll(r)
	if readErr != nil {fmt.Println(readErr)}
	var err error
	if json.Valid(fileContents) != true { 
		return err
	} // if 
	return err
} // ValidateJSON