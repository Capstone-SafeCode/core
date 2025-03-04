package A01_BrokenAccessControl

import (
	"test_capstone/src_analyser/analysis/rules/A01_BrokenAccessControl/CWE_22"
)

func RunA01Analysis(astRaw interface{}, filename string) {
	CWE_22.RunCWE22BeforeAnalysis(astRaw, filename)
}
