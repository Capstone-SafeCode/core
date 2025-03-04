package CWE_22

import (
	"fmt"
)

// Prends les données du json de la doc, RunCWE22Analysis avec les functions dangereuses puis RunCWE22Analysis avec les fonctions solutions
func RunCWE22BeforeAnalysis(astRaw interface{}, filename string) {
	isVulnerable := false

	jsonData, err := loadRules("doc/CWE-22/rule1.json")
	if err != nil {
		fmt.Printf("“Error during loading rules : %v\n", err)
		return
	}
	RunCWE22Analysis(astRaw, jsonData.Schema.DangerousFunctions, &isVulnerable, filename, true)

	if isVulnerable {
		RunCWE22Analysis(astRaw, jsonData.Schema.SafeFunctions, &isVulnerable, filename, false)
		if isVulnerable {
			var st IdentityCWE
			st.CWEId = "22"
			st.RuleId = 1
			st.Path = filename
			st.Kind = jsonData.Kind.Text
			st.ToFixIt = jsonData.ToFixIt.Text
			writeInResult(st)
			isVulnerable = false
		}
	}
}

// Extrait les noms des fonctions et leurs path avant de les envoyer à la fonction suivante
func RunCWE22Analysis(astRaw interface{}, functionsToAnalyse []interface{}, isVulnerable *bool, filename string, bootToSet bool) {
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
		if value, exists := node[key]; exists {
			return value
		}
	case []interface{}:
		for _, item := range node {
			if result := findInAST(item, key); result != nil {
				return result
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

	return false
}
