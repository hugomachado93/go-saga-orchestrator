package statemachine

import (
	"encoding/json"
	"fmt"
	"time"
)

type IStatemachineDefintion interface {
	FindNextStep(state string, event string) (string, error)
	GetName() string
	XApiKey() string
	ContextToJson() string
}

type StatemachineDefinition struct {
	Id           int64
	ClientApiKey string
	Name         string
	Context      *Statemachine
	CreatedAt    time.Time
}

func NewDefinition(clientApiKey string, context *Statemachine) (IStatemachineDefintion, error) {
	err := validateDefinition(context)
	if err != nil {
		return nil, err
	}
	return &StatemachineDefinition{ClientApiKey: clientApiKey, Name: context.Name, Context: context}, nil
}

func (stm *StatemachineDefinition) FindNextStep(state string, event string) (string, error) {
	if state == "" {
		return stm.Context.Workflow[0].NextState, nil
	}

	for _, v := range stm.Context.Workflow {

		if v.State == state && v.End {
			return "", nil
		}

		if v.State == state && v.Event == event {
			return v.NextState, nil
		}
	}
	return "", fmt.Errorf("state or event not found")
}

func (stm *StatemachineDefinition) GetName() string {
	return stm.Name
}

func (stm *StatemachineDefinition) XApiKey() string {
	return stm.ClientApiKey
}

func (stm *StatemachineDefinition) ContextToJson() string {
	val, _ := json.Marshal(stm.Context)
	return string(val)
}
