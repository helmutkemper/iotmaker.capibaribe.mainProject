package capibaribe

const KListMaxLength = 10

type servers struct {
	analytics
	Host     string  `yaml:"host"       json:"host"`
	Weight   float64 `yaml:"weight"     json:"weight"`
	OverLoad int     `yaml:"overLoad"   json:"overLoad"`
}

func (el *servers) OnExecutionStartEvent() {
	el.StartTimeCounter()
	el.OnExecutionStartCurrentExecutionsConterIncrementOne()
}

func (el *servers) OnExecutionEndWithErrorEvent() {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementErrosCounters()
	el.ResetConsecutiveSuccessCounter()
	el.SetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionTimeWithError()
}

func (el *servers) OnExecutionEndWithSuccessEvent() {
	el.OnExecutionEndCurrentExecutionsConterDecrementOne()
	el.IncrementSuccessCounters()
	el.ResetConsecutiveErrosCounter()
	el.ResetRouteHasErrorOnLastRoundFlag()
	el.AddExecutionTimeWithSuccess()
}

func NewServerStruct(host string, weight float64, overLoad int) servers {

	ret := servers{}

	ret.analytics = NewAnalytics()

	ret.Host = host
	ret.Weight = weight
	ret.OverLoad = overLoad

	return ret
}
