package genetic

type sequenceInfo struct {
	genes     string
	fitness   int
	strategy  strategyInfo
	parent    *sequenceInfo
	evolverId int
}

type strategyInfo struct {
	name         string
	start        func(strategyIndex int)
	successCount int
	results      chan *sequenceInfo
	index        int
}

type randomSource interface {
	Intn(exclusiveMax int) int
}
