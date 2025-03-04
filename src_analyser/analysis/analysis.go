package analysis

import (
	"test_capstone/src_analyser/analysis/rules/A01_BrokenAccessControl"
)

func StartAnalysis(astRaw interface{}, filename string) {
	A01_BrokenAccessControl.RunA01Analysis(astRaw, filename)
}
