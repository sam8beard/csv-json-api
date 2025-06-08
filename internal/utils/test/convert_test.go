package test

import ( 
	"testing"
	"os"
	"fmt"
)

func TestConvertToJSON(t *testing.T) { 
	invalidFilePath := "../test-files/test1-invalid.csv"
	invalidReader, err := os.Open(invalidFilePath); if err != nil {fmt.Println("Cannot open file", err)}
	invalidData, invalidErr := ConvertToJSON(invalidReader)
	_ = invalidData
	if invalidErr != nil { 
		t.Log("ConvertToJSON: Test 1 passed")
	} else { 
		t.Log("ConvertToJSON: Test 1 failed")
	} // if 
	invalidReader.Close()

	validFilePath := "../test-files/test1-valid.csv"
	validReader, err := os.Open(validFilePath); if err != nil {fmt.Println("Cannot open file", err)}
	validData, validErr := ConvertToJSON(validReader)
	_ = validData
	if validErr != nil { 
		t.Log("ConvertToJSON: Test 2 failed")
		
	} else { 
		t.Log("ConvertToJSON: Test 2 passed")
	} // if 
	validReader.Close()
	return
} // TestConvertToCSV

func TestConvertToCSV(t *testing.T) { 
	t.Log("Not yet implemented")
} // TestConvertToCSV