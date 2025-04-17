package rules

import "github.com/gin-gonic/gin"

func RunA02Analysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	RunBeforeAnalysis(resultJson, astRaw, filename, "798", "1")
}
