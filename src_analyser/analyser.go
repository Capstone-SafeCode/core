package src_analyser

import (
	"fmt"
	"log"
	"test_capstone/src_analyser/analysis"

	"github.com/gin-gonic/gin"
)

func analyseAskedCsFile(resultJson *[]gin.H, filename string) {
	return
	// astRaw, err := getCsAST(filename)

	// if err != nil {
	// 	log.Fatalf("Failed to get AST: %v\n", err)
	// }

	// //fmt.Println("Traversing JSON:")
	// //traverseJSON(astRaw, 0)

	// analysis.StartAnalysis(astRaw)
}

func analyseAskedPyFile(resultJson *[]gin.H, filename string) {
	astRaw, err := getPyAST(filename)

	if err != nil {
		log.Fatalf("Failed to get AST: %v\n", err)
	}

	analysis.StartAnalysis(resultJson, astRaw, filename)
}

func AnalyseList(listOfFiles []string) []gin.H {
	var filesList []FileManagement
	var resultJson []gin.H

	filesList = getFilesList(filesList, listOfFiles)

	analyzers := map[string]func(*[]gin.H, string){
		"py": analyseAskedPyFile,
		"cs": analyseAskedCsFile,
	}

	for _, el := range filesList {
		if analyzeFunc, ok := analyzers[el.extension]; ok {
			analyzeFunc(&resultJson, el.path)
		} else {
			fmt.Printf("Error unkown extension for '%s'\n", el.path)
		}
	}

	return resultJson
}
