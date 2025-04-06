package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func LoadJsonFromResultFile(filepath string) ([]gin.H, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error during the openning or the file : %v", err)
	}
	defer file.Close()

	var data []gin.H
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("error during json's decodage : %v", err)
	}

	return data, nil
}
