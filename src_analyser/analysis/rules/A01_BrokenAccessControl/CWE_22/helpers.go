package CWE_22

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type IdentityCWE struct {
	CWEId   string
	RuleId  int
	Path    string
	Line    int
	Kind    string
	ToFixIt string
}

type SecurityRules struct {
	Schema struct {
		DangerousFunctions []interface{} `json:"dangerous_functions"`
		SafeFunctions      []interface{} `json:"safe_functions"`
	} `json:"schema"`
	ToFixIt struct {
		Text string `json:"text"`
	} `json:"ToFixIt"`
	Kind struct {
		Text string `json:"text"`
	} `json:"Kind"`
}

func SplitStringByDot(input string) []string {
	return strings.Split(input, ".")
}

func writeInResult(st IdentityCWE) {
	filename := "result.json"

	var data []gin.H

	file, err := os.Open(filename)
	if err == nil {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			fmt.Printf("Erreur lors du décodage du JSON existant : %v\n", err)
		}
		file.Close()
	} else if !os.IsNotExist(err) {
		fmt.Printf("Erreur lors de l'ouverture du fichier : %v\n", err)
		return
	}

	newEntry := gin.H{
		"CWE":     st.CWEId,
		"RuleId":  st.RuleId,
		"Path":    st.Path,
		"Kind":    st.Kind,
		"ToFixIt": st.ToFixIt,
	}
	data = append(data, newEntry)

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Erreur lors de l'ouverture du fichier pour écriture : %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("Erreur lors de l'écriture dans le fichier : %v\n", err)
	}
}

func loadRules(filepath string) (SecurityRules, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return SecurityRules{}, err
	}
	defer file.Close()

	var rules SecurityRules
	if err := json.NewDecoder(file).Decode(&rules); err != nil {
		return SecurityRules{}, err
	}

	return rules, nil
}
