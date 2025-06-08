package test

import ( 
	"testing"
	"os"
	"fmt"
)

func TestValidateCSV(t *testing.T) { 
	t.Log("Not yet implemented")
} // TestValidateCSV

func TestValidateJSON(t *testing.T) { 
	invalidFilePath := "../test-files/test2-invalid.json"
	invalidReader, err := os.Open(invalidFilePath); if err != nil {fmt.Println("Cannot open file")}
	invalidErr := ValidateJSON(invalidReader) 
	if invalidErr != nil { 
		t.Log("Test 1 passed")
		
	} else { 
		t.Log("Test 1 failed")
	} // if 
	invalidReader.Close()

	validFilePath := "../test-files/test2-valid.json"
	validReader, err := os.Open(validFilePath); if err != nil {fmt.Println("Cannot open file")}
	validErr := ValidateJSON(validReader)
	if validErr != nil { 
		t.Log("Test 2 passed")
		
	} else { 
		t.Log("Test 2 failed")
	} // if 
	validReader.Close()
	return

} // TestValidateJSON
