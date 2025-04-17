package rules

import "github.com/gin-gonic/gin"

func RunA05Analysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	RunBeforeAnalysis(resultJson, astRaw, filename, "613", "1")
}
