package test

import ( 
	"io"
	"encoding/json"
	"errors"
	"io/ioutil"
	"fmt"
)
func ValidateCSV(r io.Reader) error { 
	// Check if suffix .csv
	return nil
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
	if readErr != nil {fmt.Println("Can't read")}
	var err error
	if json.Valid(fileContents) != true { 
		err = errors.New("Invalid file format")
		return err
	} // if 
	return err
} // ValidateJSON