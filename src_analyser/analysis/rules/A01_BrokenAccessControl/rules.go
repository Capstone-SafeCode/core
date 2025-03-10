package A01_BrokenAccessControl

func RunA01Analysis(astRaw interface{}, filename string) {
	RunBeforeAnalysis(astRaw, filename, "22", "1")

	RunBeforeAnalysis(astRaw, filename, "200", "1")
}
