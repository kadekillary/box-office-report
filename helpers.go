package main

import (
	"strings"
)

// ExtractFileName out of URL and convert file format
func ExtractFileName(url string) string {
	// TODO: Could refactor to write each week's results into a separate file
	splitURL := strings.Split(url, "/")
	extractedFileName := splitURL[len(splitURL)-1]
	fileName := strings.Replace(extractedFileName, "html", "csv", -1)
	return fileName
}
