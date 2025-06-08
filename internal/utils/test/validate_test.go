package test

import ( 
	"testing"
	"os"
	"fmt"
)

func TestValidateCSV(t *testing.T) { 
	// t.Log("Not yet implemented")
	invalidFilePath := "../test-files/test1-invalid.csv"
	invalidReader, err := os.Open(invalidFilePath); if err != nil {fmt.Println("Cannot open file", err)}
	invalidErr := ValidateCSV(invalidReader)
	if invalidErr != nil { 
		t.Log("ValidateCSV: Test 1 passed")
	} else { 
		t.Log("ValidateCSV: Test 1 failed")
	} // if 
	invalidReader.Close()

	validFilePath := "../test-files/test1-valid.csv"
	validReader, err := os.Open(validFilePath); if err != nil {fmt.Println("Cannot open file")}
	validErr := ValidateCSV(validReader)
	if validErr != nil { 
		t.Log("ValidateCSV: Test 2 failed")
		
	} else { 
		t.Log("ValidateCSV: Test 2 passed")
	} // if 
	validReader.Close()
	return

} // TestValidateCSV

func TestValidateJSON(t *testing.T) { 
	invalidFilePath := "../test-files/test2-invalid.json"
	invalidReader, err := os.Open(invalidFilePath); if err != nil {fmt.Println("Cannot open file")}
	invalidErr := ValidateJSON(invalidReader) 
	if invalidErr != nil { 
		t.Log("ValidateJSON: Test 1 failed")
		
	} else { 
		t.Log("ValidateJSON: Test 1 passed")
	} // if 
	invalidReader.Close()

	validFilePath := "../test-files/test2-valid.json"
	validReader, err := os.Open(validFilePath); if err != nil {fmt.Println("Cannot open file")}
	validErr := ValidateJSON(validReader)
	if validErr != nil { 
		t.Log("ValidateJSON: Test 2 failed")
		
	} else { 
		t.Log("ValidateJSON: Test 2 passed")
	} // if 
	validReader.Close()
	return

} // TestValidateJSON
