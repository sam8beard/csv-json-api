package utils

import ( 
	"fmt"
	"os"
)

func DownloadFile(url string) (io.ReadCloser) { 
	var file *os.File
	_ = file
	fmt.Println("Testing")
	return file
} // DownloadFile