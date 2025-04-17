package rules

import "github.com/gin-gonic/gin"

func RunA04Analysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	RunBeforeAnalysis(resultJson, astRaw, filename, "352", "1")
}
