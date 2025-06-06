# CSV-JSON Converter API
 A simple, stateless backend API written in Go for converting `.csv` files to `.json` files and vice versa. This API leverages Go's concurrency features to perform fast and efficient conversions. Uploaded files are processed and returned in a `.zip` archive. I've been eager to learn Go for some time, and this seemed like a suitable first project to tackle while familiarizing myself with some of the basics.

## Features
- Accepts local file uploads (`multipart/form-data`)
- Accepts file URLs to download and convert
- Converts between `.csv` and `.json`
- Returns a `.zip` file of converted results
- Simple CLI testing with `curl`
### Planned features (stretch)
- Stream-based processing for large files 

## Installation

`go build -o api main.go`

## Usage / Examples 
```bash
curl -X POST http://localhost:8080/convert \
  -F "files=@path/to/file.csv" \
  -F "urls=https://example.com/data.json" \
  --output result.zip``
```
