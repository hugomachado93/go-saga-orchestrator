package saga

import (
	"time"
)

type SagaStatus string

type ResponseStatus string

type Response struct {
	SagaName string `json:"name"`
	SagaUUID string `json:"uuid"`
	Payload  string `json:"payload"`
	Event    string `json:"event"`
}

type Saga struct {
	Id           *int64    `db:"id"`
	ApiKey       string    `db:"api_key"`
	UUID         string    `db:"uuid"`
	Name         string    `db:"name"`
	State        string    `db:"state"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	LastUpdate   time.Time `db:"last_update"`
	SagaCommands []SagaCommand
}

func NewSaga(apiKey string, uuid string, name string) *Saga {
	s := &Saga{ApiKey: apiKey, Name: name, UUID: uuid, SagaCommands: make([]SagaCommand, 0)}
	return s
}

func (s *Saga) PrepareNextCommand(payload string, nextState string, delay *int) {

	s.LastUpdate = time.Now()
	s.State = nextState
	s.Status = "EXECUTING"
	s.SagaCommands = append(s.SagaCommands, NewCommand(s.State, s.Id, nil))
	s.LastUpdate = time.Now()
	s.CreatedAt = time.Now()
}

func (s *Saga) FinishSaga() {
	s.Status = "COMPLETED"
	s.LastUpdate = time.Now()
}

func (s *Saga) RequestCommand() {
	if len(s.SagaCommands) > 0 {
		s.SagaCommands[len(s.SagaCommands)-1].RequestCommand()
	}
}

func (s *Saga) ReceiveResponse() {
	if len(s.SagaCommands) > 0 {
		s.SagaCommands[len(s.SagaCommands)-1].ReceiveResponse()
	}
}

func (s *Saga) HasSagaTimedout() bool {
	if len(s.SagaCommands) > 0 {
		return s.SagaCommands[len(s.SagaCommands)-1].HasCommandTimeout()
	}
	return false
}
