package rules

import "github.com/gin-gonic/gin"

func RunA01Analysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	RunBeforeAnalysis(resultJson, astRaw, filename, "22", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "200", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "201", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "285", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "287", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "384", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "639", "1")
}
