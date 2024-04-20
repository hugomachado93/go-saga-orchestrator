package requests

import (
	"main/internal/domain/statemachine"
)

type Stmc struct {
	State     string `json:"state"`
	Event     string `json:"event"`
	NextState string `json:"nextState"`
	Delay     *int   `json:"delay"`
	End       bool   `json:"end"`
}

type Statemachine struct {
	Name     string `json:"name"`
	Workflow []Stmc `json:"workflow"`
}

func (stm *Statemachine) ToStateMachineSeetings() *statemachine.Statemachine {
	stmcs := make([]statemachine.Stmc, 0)
	for _, v := range stm.Workflow {
		stmc := statemachine.Stmc{State: v.State, Event: v.Event, NextState: v.NextState, Delay: v.Delay, End: v.End}
		stmcs = append(stmcs, stmc)
	}
	return &statemachine.Statemachine{Name: stm.Name, Workflow: stmcs}
}
