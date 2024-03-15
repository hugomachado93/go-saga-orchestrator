package saga

type nextRollback interface {
	NextRollback(step Command) and
	And() nextStep
	Build()
}

type nextStep interface {
	NextStep(step Command) nextRollback
}

type and interface {
	And() nextStep
	Build()
}

var sagaMap map[string]SagaList = make(map[string]SagaList)

type SagaList [][]Command

type SagaBuilder struct {
	sagaName     string
	nextStep     Command
	rollbackStep Command
	saga         SagaList
}

func NewSagaDefinition(sagaName string) nextStep {
	s := &SagaBuilder{sagaName: sagaName, saga: make([][]Command, 0)}
	return s
}

func (s *SagaBuilder) NextStep(step Command) nextRollback {
	s.nextStep = step
	return s
}

func (s *SagaBuilder) NextRollback(step Command) and {
	s.rollbackStep = step
	return s
}

func (s *SagaBuilder) And() nextStep {
	s.insert()
	return s
}

func (s *SagaBuilder) Build() {
	s.insert()
	sagaMap[s.sagaName] = s.saga
}

func (s *SagaBuilder) insert() {
	arr := make([]Command, 0)
	arr = append(arr, s.nextStep)
	arr = append(arr, s.rollbackStep)
	s.saga = append(s.saga, arr)
}

func GetNextCommanfWithIndex(sagaName string, idx int, isRollback bool) Command {
	size := len(sagaMap[sagaName])
	if idx >= size || idx < 0 {
		return ""
	}
	return sagaMap[sagaName][idx][getIdy(isRollback)]
}

func getIdy(isRollback bool) int {
	if isRollback {
		return 1
	}
	return 0
}
