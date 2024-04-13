package saga

import "fmt"

type sagaStateMachine [][]string

var sagaMap2 map[string][][]string = make(map[string][][]string)

type SagaBuilder2 struct {
	sagaName     string
	state        string
	event        string
	nextState    string
	stateMachine [][]string
}

func NewSagaBuilder2(sagaName string) *SagaBuilder2 {
	return &SagaBuilder2{sagaName: sagaName, stateMachine: make([][]string, 0)}
}

func (s *SagaBuilder2) OnState(state string) *SagaBuilder2 {
	s.state = state
	return s
}

func (s *SagaBuilder2) WithEvent(event string) *SagaBuilder2 {
	s.event = event
	return s
}

func (s *SagaBuilder2) NextState(nextState string) *SagaBuilder2 {
	s.nextState = nextState
	return s
}

func (s *SagaBuilder2) And() *SagaBuilder2 {
	s.insert()
	return s
}

func (s *SagaBuilder2) Build() {
	s.insert()
	sagaMap2[s.sagaName] = s.stateMachine
}

func (s *SagaBuilder2) insert() {
	arr := make([]string, 0)
	arr = append(arr, s.state)
	arr = append(arr, s.event)
	arr = append(arr, s.nextState)
	s.stateMachine = append(s.stateMachine, arr)
}

func getNextStep(sagaName string, currentSate string, event string) (string, error) {
	sagalst := sagaMap2[sagaName]
	for _, vlst := range sagalst {
		if vlst[0] == currentSate && vlst[1] == event {
			return vlst[2], nil
		}
	}

	return "", fmt.Errorf("asdsad")
}
