package analysis

import (
	"test_capstone/src_analyser/analysis/rules"

	"github.com/gin-gonic/gin"
)

func StartPyAnalysis(resultJson *[]gin.H, astRaw interface{}, filename string) {
	rules.RunA01Analysis(resultJson, astRaw, filename)
	rules.RunA02Analysis(resultJson, astRaw, filename)
	rules.RunA03Analysis(resultJson, astRaw, filename)
	rules.RunA04Analysis(resultJson, astRaw, filename)
	rules.RunA05Analysis(resultJson, astRaw, filename)
}
