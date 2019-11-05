package capibaribe

const KExecutionTimeListLength = 10

type servers struct {
	executionTimeMax     int64
	executionTimeMin     int64
	executionTimeList    []int64
	executionTimeAverage int64
	consecutiveErrors    int
	consecutiveSuccess   int
	errors               int
	success              int
	lastRoundError       bool
	Host                 string  `yaml:"host"       json:"host"`
	Weight               float64 `yaml:"weight"     json:"weight"`
	OverLoad             int     `yaml:"overLoad"   json:"overLoad"`
}

func (el *servers) AddExecutionTime(duration int64) {

	if duration > el.executionTimeMax {
		el.executionTimeMax = duration
	}

	if el.executionTimeMin == 0 {
		el.executionTimeMin = duration
	} else if el.executionTimeMin > duration {
		el.executionTimeMin = duration
	}

	el.executionTimeList = append(el.executionTimeList, duration)
	if len(el.executionTimeList) > KExecutionTimeListLength {
		el.executionTimeList = el.executionTimeList[1:]
	}

	el.executionTimeAverage = 0
	for _, value := range el.executionTimeList {
		el.executionTimeAverage += value
	}

	el.executionTimeAverage = el.executionTimeAverage / int64(len(el.executionTimeList))
}

func NewServerStruct() servers {
	ret := servers{}
	ret.executionTimeList = make([]int64, 0)

	return ret
}
