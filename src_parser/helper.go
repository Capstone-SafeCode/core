package src_parser

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func contains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}

func getExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}

func isItFile(filepath string) (int, string) {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("Error path '%s' doesn't exist\n", filepath)
		return -1, ""
	}
	if err != nil {
		fmt.Printf("Error during path verification '%s': %v\n", filepath, err)
		return -1, ""
	}

	if info.IsDir() {
		// Folder
		return 0, ""
	} else {
		// File
		return 1, getExtension(filepath)
	}
}

func findItsIndex(slice []string, filepath string) int {
	sort.Strings(slice)
	index := sort.SearchStrings(slice, filepath)

	if index <= len(slice) && slice[index] == filepath {
		return index
	} else {
		return -1
	}
}
