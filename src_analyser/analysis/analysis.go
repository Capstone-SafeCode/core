package analysis

import (
	"test_capstone/src_analyser/analysis/rules/A01_BrokenAccessControl"

	"github.com/gin-gonic/gin"
)

func StartAnalysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	A01_BrokenAccessControl.RunA01Analysis(resultJson, astRaw, filename)
}
