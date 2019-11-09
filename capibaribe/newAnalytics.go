package capibaribe

// pt_br: Prepara o struct do tipo AnalyticsCheck para ser usado
// en: Mounts a struct type AnalyticsCheck for a correct use
func NewAnalytics() Analytics {
	ret := Analytics{}
	ret.ExecutionDurationEntireList = make([]executionInfo, 0)
	ret.ExecutionDurationSuccessList = make([]executionInfo, 0)
	ret.ExecutionDurationErrorList = make([]executionInfo, 0)

	return ret
}
