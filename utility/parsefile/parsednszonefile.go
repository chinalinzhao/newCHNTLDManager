package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	// Read the DNS zone file
	filePath := "/c:/Users/polo/OneDrive/work/project/毛总/newChn/newCHNTLDManager/utility/parsefile/parsednszonefile.go"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Split the file content into lines
	lines := strings.Split(string(content), "\n")

	// Process each line
	for _, line := range lines {
		// Skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, ";") {
			continue
		}

		// Parse the DNS record
		// TODO: Implement your parsing logic here

		// Print the parsed record
		fmt.Println(line)
	}
}
