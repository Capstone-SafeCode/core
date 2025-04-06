package A01_BrokenAccessControl

import "github.com/gin-gonic/gin"

func RunA01Analysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	RunBeforeAnalysis(resultJson, astRaw, filename, "22", "1")
	RunBeforeAnalysis(resultJson, astRaw, filename, "200", "1")
}
