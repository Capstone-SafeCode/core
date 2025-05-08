package rules

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Prends les données du json de la doc, RunCWE22Analysis avec les functions dangereuses puis RunCWE22Analysis avec les fonctions solutions
func RunBeforeAnalysis(resultJson *[]gin.H, astRaw interface{}, filename string, whichCWE string, whichRule string) {
	isVulnerable := false
	var whichLines []int

	jsonData, err := loadRules("doc/CWE-" + whichCWE + "/rule" + whichRule + ".json")
	if err != nil {
		fmt.Printf("“Error during loading rules : %v\n", err)
		return
	}
	RunAnalysis(astRaw, jsonData.Schema.DangerousFunctions, &isVulnerable, &whichLines, filename, true)

	if isVulnerable {
		RunAnalysis(astRaw, jsonData.Schema.SafeFunctions, &isVulnerable, &whichLines, filename, false)
		if isVulnerable {
			*resultJson = append(*resultJson, gin.H{
				"CWE":     whichCWE,
				"RuleId":  whichRule,
				"Path":    filename,
				"Lines":   whichLines,
				"Kind":    jsonData.Kind.Text,
				"ToFixIt": jsonData.ToFixIt.Text,
			})
			isVulnerable = false
		}
	}
}

// Extrait les noms des fonctions et leurs path avant de les envoyer à la fonction suivante
func RunAnalysis(astRaw interface{}, functionsToAnalyse []interface{}, isVulnerable *bool, whichLines *[]int, filename string, boolToSet bool) {
	nbToFind := 0
	nbFound := 0

	for _, dFunc := range functionsToAnalyse {
		dFuncList, ok := dFunc.([]interface{})
		if !ok || len(dFuncList) < 2 {
			fmt.Println("Error: incorrect format in DangerousFunctions/SafeFunctions")
			continue
		}

		functionName, ok := dFuncList[0].(string)
		if !ok {
			fmt.Println("Error: functionName is not a string")
			continue
		}

		dFunctionPathList, ok := dFuncList[1].([]interface{})
		if !ok {
			fmt.Println("Error: dFunctionPath is not a list")
			continue
		}

		var dFunctionPath []string
		for _, path := range dFunctionPathList {
			if strPath, ok := path.(string); ok {
				dFunctionPath = append(dFunctionPath, strPath)
			}
		}

		nbToFind += len(SplitStringByDot(functionName))

		nbFound += analyzeAST(astRaw, functionName, dFunctionPath, whichLines)
	}

	if nbFound == nbToFind {
		*isVulnerable = boolToSet
	}
}

// Boucle sur les path et, si vrai, set isVulnerable au booléen demandé
func analyzeAST(astRaw interface{}, functionName string, dFunctionPath []string, whichLines *[]int) int {
	tempNb := 0

	functionParts := SplitStringByDot(functionName)

	for _, path := range dFunctionPath {
		pathParts := SplitStringByDot(path)
		if searchInAST(astRaw, functionParts, pathParts, whichLines) {
			tempNb++
		}
	}
	return tempNb
}

// Boucle sur tous les body et applique le path demandé
func searchInAST(ast interface{}, functionParts []string, pathParts []string, whichLines *[]int) bool {
	if explorePath(ast, pathParts, functionParts, whichLines) {
		return true
	}

	switch node := ast.(type) {
	case map[string]interface{}:
		for _, value := range node {
			if searchInAST(value, functionParts, pathParts, whichLines) {
				return true
			}
		}
	case []interface{}:
		for _, item := range node {
			if searchInAST(item, functionParts, pathParts, whichLines) {
				return true
			}
		}
	}

	return false
}

// Utilise les deux fonctions en dessous pour vérifier les données qu'on lui donne
func explorePath(ast interface{}, pathParts []string, functionParts []string, whichLines *[]int) bool {
	currentNode := ast
	var parentNode interface{} = nil

	for _, key := range pathParts {
		nextNode := findInAST(currentNode, key)
		if nextNode == nil {
			return false
		}

		parentNode = currentNode
		currentNode = nextNode
	}

	if matchFunction(currentNode, functionParts) {
		if nodeMap, ok := parentNode.(map[string]interface{}); ok {
			if line, ok := nodeMap["lineno"].(float64); ok {
				lineNum := int(line)
				if !ContainsLine(*whichLines, lineNum) {
					*whichLines = append(*whichLines, lineNum)
				}
			}
		}
		return true
	}

	return false
}

// Cherche un élément dans l'interface (l'AST)
func findInAST(ast interface{}, key string) interface{} {
	switch node := ast.(type) {
	case map[string]interface{}:
		// Accès classique à une clé
		if value, exists := node[key]; exists {
			return value
		}
	case []interface{}:
		// Si la clé est un index (ex: "0"), on l'interprète comme tel
		if idx, err := strconv.Atoi(key); err == nil {
			if idx >= 0 && idx < len(node) {
				return node[idx]
			}
		}
	}

	return nil
}

// Vérifie si la valeur trouvée est la même que demandée
func matchFunction(ast interface{}, functionParts []string) bool {
	if strValue, ok := ast.(string); ok {
		for _, expected := range functionParts {
			if strValue == expected {
				return true
			}
		}
	}

	if nodeMap, ok := ast.(map[string]interface{}); ok {
		if idVal, ok := nodeMap["id"]; ok {
			if strID, ok := idVal.(string); ok {
				for _, expected := range functionParts {
					if strID == expected {
						return true
					}
				}
			}
		}

		if typeVal, ok := nodeMap["_type"]; ok {
			if strType, ok := typeVal.(string); ok {
				for _, expected := range functionParts {
					if strType == expected {
						return true
					}
				}
			}
		}
	}

	return false
}
