package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func help() {
	fmt.Println("usage: ./parser -path=[f]")
	fmt.Println("Options corresponding:")
	fmt.Println("f      : Folder's path from which the analysis will start (ex : src/)")
}

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

func writeInTXT(text string) {
	filename := "to_analyse.txt"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error in openning : %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		fmt.Printf("Error in writting : %v\n", err)
		return
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
