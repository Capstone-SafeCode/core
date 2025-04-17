package rules

import "github.com/gin-gonic/gin"

func RunA03Analysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	RunBeforeAnalysis(resultJson, astRaw, filename, "20", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "74", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "79", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "89", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "94", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "917", "1")
}
