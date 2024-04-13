package statemachine

import (
	"main/internal/err"
	"time"
)

type Stmc struct {
	State     string `json:"state"`
	Event     string `json:"event"`
	NextState string `json:"nextState"`
}

type Statemachine struct {
	Name     string `json:"name"`
	Workflow []Stmc `json:"workflow"`
}

type StatemachineSettings struct {
	Id           int64
	ClientApiKey string
	Name         string
	Context      Statemachine
	CreatedAt    time.Time
}

func (stm *StatemachineSettings) FindNextStep(state string, event string) (string, error) {
	if state == "" {
		return stm.Context.Workflow[0].NextState, nil
	}

	for _, v := range stm.Context.Workflow {
		if v.State == state && v.Event == event {
			return v.NextState, nil
		}
	}
	return "", err.ErroStateOrEventNotFound
}
