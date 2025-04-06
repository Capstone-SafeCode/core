package src_analyser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func createEmptyJSON(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("[]")
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture dans le fichier : %v", err)
	}

	return nil
}

type FileManagement struct {
	path      string
	extension string
}

func getPyAST(filename string) (interface{}, error) {
	cmd := exec.Command("python3", "ast/myast.py", filename)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("command failed with error: %v, stderr: %s", err, stderr.String())
	}

	//fmt.Println("Raw from Python script:\n", out.String())

	var raw interface{}
	if err := json.Unmarshal(out.Bytes(), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return raw, nil
}

func getFilesList(filesList []FileManagement, listOfFiles []string) []FileManagement {
	for _, line := range listOfFiles {
		split := strings.Fields(line)
		if len(split) < 2 {
			continue
		}

		tempStruc := FileManagement{
			path:      split[0],
			extension: split[1],
		}

		filesList = append(filesList, tempStruc)
	}

	return filesList
}
