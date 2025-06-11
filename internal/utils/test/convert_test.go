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

	validFilePath2 := "../test-files/test4-valid.csv"
	validReader2, err := os.Open(validFilePath2); if err != nil {fmt.Println("Cannot open file", err)}
	validData2, validErr2 := ConvertToJSON(validReader2)
	_ = validData2
	if validErr2 != nil { 
		t.Log("ConvertToJSON: Test 3 failed")
		
	} else { 
		t.Log("ConvertToJSON: Test 3 passed")
	} // if 
	validReader2.Close()

	return
} // TestConvertToCSV

func TestConvertToCSV(t *testing.T) { 
	invalidFilePath := "../test-files/test2-invalid.json"
	invalidReader, err := os.Open(invalidFilePath); if err != nil {fmt.Println("Cannot open file", err)}
	invalidData, invalidErr := ConvertToCSV(invalidReader)
	_ = invalidData
	if invalidErr != nil { 
		t.Log("ConvertToCSV: Test 1 passed")
	} else { 
		t.Log("ConvertToCSV: Test 1 failed")
	} // if 
	invalidReader.Close()

	validFilePath := "../test-files/test2-valid.json"
	validReader, err := os.Open(validFilePath); if err != nil {fmt.Println("Cannot open file", err)}
	validData, validErr := ConvertToCSV(validReader)
	_ = validData
	if validErr != nil { 
		t.Log("ConvertToCSV: Test 2 failed")
		
	} else { 
		t.Log("ConvertToCSV: Test 2 passed")
	} // if 
	validReader.Close()

	// invalidFilePath := "../test-files/test2-invalid.json"
	// invalidReader, err := os.Open(invalidFilePath); if err != nil {fmt.Println("Cannot open file", err)}
	// invalidData, invalidErr := ConvertToCSV(invalidReader)
	// _ = invalidData
	// if invalidErr != nil { 
	// 	t.Log("ConvertToCSV: Test 1 passed")
	// } else { 
	// 	t.Log("ConvertToCSV: Test 1 failed")
	// } // if 
	// invalidReader.Close()

	validFilePath2 := "../test-files/test3-valid.json"
	validReader2, err := os.Open(validFilePath2); if err != nil {fmt.Println("Cannot open file", err)}
	validData2, validErr2 := ConvertToCSV(validReader2)
	_ = validData2
	if validErr2 != nil { 
		t.Log("ConvertToCSV: Test 3 failed")
		
	} else { 
		t.Log("ConvertToCSV: Test 3 passed")
	} // if 
	validReader2.Close()
	return
} // TestConvertToCSV