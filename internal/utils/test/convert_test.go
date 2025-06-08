package test

import ( 
	"testing"
	"os"
	"fmt"
)

func TestConvertToJSON(t *testing.T) { 
	t.Log("Not yet implemented")
	validFilePath := "../test-files/test1-valid.csv"
	validReader, err := os.Open(validFilePath); if err != nil {fmt.Println("Cannot open file")}
	file, err := ConvertToJSON(validReader)
	_ = file
} // TestConvertToCSV

func TestConvertToCSV(t *testing.T) { 
	t.Log("Not yet implemented")
} // TestConvertToCSV