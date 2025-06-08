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
	validErr := ValidateJSON(invalidReader) 
	if validErr != nil { 
		t.Log("Test passed")
		return
	} 
	t.Log("Test failed")
	return

} // TestValidateJSON
