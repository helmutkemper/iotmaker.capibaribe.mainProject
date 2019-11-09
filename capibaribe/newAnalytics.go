package capibaribe

// pt_br: Prepara o struct do tipo Analytics para ser usado
// en: Mounts a struct type Analytics for a correct use
func NewAnalytics() Analytics {
	ret := Analytics{}
	ret.ExecutionDurationEntireList = make([]executionInfo, 0)
	ret.ExecutionDurationSuccessList = make([]executionInfo, 0)
	ret.ExecutionDurationErrorList = make([]executionInfo, 0)

	return ret
}
