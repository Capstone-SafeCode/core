package A01_BrokenAccessControl

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Prends les données du json de la doc, RunCWE22Analysis avec les functions dangereuses puis RunCWE22Analysis avec les fonctions solutions
func RunBeforeAnalysis(resultJson *[]gin.H, astRaw interface{}, filename string, whichCWE string, whichRule string) {
	isVulnerable := false

	jsonData, err := loadRules("doc/CWE-" + whichCWE + "/rule" + whichRule + ".json")
	if err != nil {
		fmt.Printf("“Error during loading rules : %v\n", err)
		return
	}
	RunAnalysis(astRaw, jsonData.Schema.DangerousFunctions, &isVulnerable, filename, true)

	if isVulnerable {
		RunAnalysis(astRaw, jsonData.Schema.SafeFunctions, &isVulnerable, filename, false)
		if isVulnerable {
			*resultJson = append(*resultJson, gin.H{
				"CWE":     whichCWE,
				"RuleId":  whichRule,
				"Path":    filename,
				"Kind":    jsonData.Kind.Text,
				"ToFixIt": jsonData.ToFixIt.Text,
			})
			isVulnerable = false
		}
	}
}

// Extrait les noms des fonctions et leurs path avant de les envoyer à la fonction suivante
func RunAnalysis(astRaw interface{}, functionsToAnalyse []interface{}, isVulnerable *bool, filename string, bootToSet bool) {
	for _, dFunc := range functionsToAnalyse {
		dFuncList, ok := dFunc.([]interface{})
		if !ok || len(dFuncList) < 2 {
			fmt.Println("Error: incorrect format in DangerousFunctions")
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

		analyzeAST(astRaw, functionName, dFunctionPath, isVulnerable, bootToSet)
	}
}

// Boucle sur les path et, si vrai, set isVulnerable au booléen demandé
func analyzeAST(astRaw interface{}, functionName string, dFunctionPath []string, isVulnerable *bool, bootToSet bool) {
	functionParts := SplitStringByDot(functionName)

	for _, path := range dFunctionPath {
		pathParts := SplitStringByDot(path)
		if searchInAST(astRaw, functionParts, pathParts) {
			*isVulnerable = bootToSet
			return
		}
	}
}

// Boucle sur tous les body et applique le path demandé
func searchInAST(ast interface{}, functionParts []string, pathParts []string) bool {
	if explorePath(ast, pathParts, functionParts) {
		return true
	}

	switch node := ast.(type) {
	case map[string]interface{}:
		for key, value := range node {
			if key == "body" {
				if bodyArray, ok := value.([]interface{}); ok {
					for _, subAST := range bodyArray {
						if searchInAST(subAST, functionParts, pathParts) {
							return true
						}
					}
				}
			} else {
				if searchInAST(value, functionParts, pathParts) {
					return true
				}
			}
		}
	case []interface{}:
		for _, item := range node {
			if searchInAST(item, functionParts, pathParts) {
				return true
			}
		}
	}

	return false
}

// Utilise les deux fonctions en dessous pour vérifier les données qu'on lui donne
func explorePath(ast interface{}, pathParts []string, functionParts []string) bool {
	currentNode := ast

	for _, key := range pathParts {
		nextNode := findInAST(currentNode, key)

		if nextNode == nil {
			return false
		}

		currentNode = nextNode
	}

	return matchFunction(currentNode, functionParts)
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

	// Cas spécial : si c’est un map et qu’il a une clé "id" contenant un des noms attendus
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
	}

	return false
}
