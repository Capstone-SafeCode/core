package rules

import (
	"encoding/json"
	"os"
	"strings"
)

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

func ContainsLine(lines []int, line int) bool {
	for _, l := range lines {
		if l == line {
			return true
		}
	}
	return false
}
