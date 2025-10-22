# CSV-JSON Converter API
 A simple, stateless backend API written in Go for converting `.csv` files to `.json` files and vice versa. This API leverages Go's concurrency features to perform fast and efficient conversions. Uploaded files are processed and returned in a `.zip` archive. I've been eager to learn Go for some time, and this seemed like a suitable first project to tackle while familiarizing myself with some of the basics.

## NOTES
 - This API only supports flat JSON objects for conversion purposes. 
 - All CSV files must follow a standard structure (see below for CSV format expectations).

## CSV Format Requirements 
To ensure reliable conversion, .csv files must conform to the following rules:
-	Must use commas (,) as delimiters (no tabs, semicolons, etc.)
-	Must include a header row as the first line (used as JSON object keys)
-	Each row must contain the same number of fields as the header
- Fields may be optionally wrapped in double quotes ("), especially if they contain commas
-	Quoted fields must follow standard escaping conventions (e.g., He said ""hello"")
- UTF-8 encoded files are supported; no BOM (Byte Order Mark)

In short, the CSV format generally adheres to RFC 4180.

## Features
- Accepts local file uploads (`multipart/form-data`)
- Accepts file URLs to download and convert
- Converts between `.csv` and `.json`
- Returns a `.zip` file of converted results
- Includes a `log.json` file detailing which files were skipped or converted

### Planned features (stretch)
- Stream-based processing for large files 

## How to build and run

`./build.sh`

`./api`

or to run in the background

`./api &`

## Usage / Examples 
```bash
curl -X POST http://localhost:8080/convert \
  -F "files=@path/to/file.csv" \
  -F "urls=https://example.com/data.json" \
  --output result.zip``
```

The response will contain 
- Converted `.json` or `.csv` files 
- A `log.json` file summarizing which files were processed successfully and which were skipped (with reasons)