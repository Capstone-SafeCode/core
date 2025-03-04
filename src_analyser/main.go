package main

import (
	"fmt"
	"log"
	"test_capstone/src_analyser/analysis"
)

func analyseAskedCsFile(filename string) {
	return
	// astRaw, err := getCsAST(filename)

	// if err != nil {
	// 	log.Fatalf("Failed to get AST: %v\n", err)
	// }

	// //fmt.Println("Traversing JSON:")
	// //traverseJSON(astRaw, 0)

	// analysis.StartAnalysis(astRaw)
}

func analyseAskedPyFile(filename string) {
	astRaw, err := getPyAST(filename)

	if err != nil {
		log.Fatalf("Failed to get AST: %v\n", err)
	}

	analysis.StartAnalysis(astRaw, filename)
}

func main() {
	var filesList []FileManagement

	filesList = getFilesList(filesList)

	analyzers := map[string]func(string){
		"py": analyseAskedPyFile,
		"cs": analyseAskedCsFile,
	}

	err := createEmptyJSON("result.json")
	if err != nil {
		fmt.Println("Erreur :", err)
	}

	for _, el := range filesList {
		if analyzeFunc, ok := analyzers[el.extension]; ok {
			analyzeFunc(el.path)
		} else {
			fmt.Printf("Error unkown extension for '%s'\n", el.path)
		}
	}
}
