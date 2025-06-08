package utils 

import ( 
	"io"
	"encoding/json"
	"errors"
)
func ValidateCSV(r io.Reader) error { 
	// Check if suffix .csv
	
} // ValidateCSV

/*
I know that encoding/json literally has a built in function to 
validate json objects and this is redundant in every way, 
but I'm instantiating my own function for my sanity.

I'm just wrapping the call to json.Validate in this function 
and returning an error instead of bool.
*/
func ValidateJSON(r io.Reader) error {
	fileContents := r.ReadAll() 
	err := nil 
	if json.Valid(fileContents) != true { 
		error := errors.New("Invalid file format")
	} // if 
	return err
} // ValidateJSON