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

type Status string

const (
	WAITING_PROCESS   Status = "WAITING_PROCESS"
	REQUEST_SENT      Status = "REQUEST_SENT"
	RESPONSE_RECEIVED Status = "RESPONSE_RECEIVED"
)

type Saga struct {
	Id           int64
	ApiKey       string
	SagaName     string
	SagaUUID     string
	Payload      string
	CurrentState string
	Status       Status
	Timeout      *time.Time
	CreatedAt    time.Time
	LastUpdate   time.Time
}

func NewSaga(apiKey string, uuid string, name string) *Saga {
	return &Saga{ApiKey: apiKey, SagaUUID: uuid, SagaName: name, CreatedAt: time.Now()}
}

func (s *Saga) PrepareNextCommand(payload string, nextState string) {
	s.CurrentState = nextState
	s.Payload = payload
	s.Status = WAITING_PROCESS
	s.LastUpdate = time.Now()
	s.CreatedAt = time.Now()
}

func (s *Saga) RequestSaga() {
	s.Status = REQUEST_SENT
	s.LastUpdate = time.Now()
	t := time.Now().Add(time.Duration(time.Duration.Seconds(120)))
	s.Timeout = &t
}

func (s *Saga) ReceiveResponse() {
	if s.Status == REQUEST_SENT {
		s.Status = RESPONSE_RECEIVED
	}
}

func (s *Saga) HasSagaTimedout() bool {
	tn := time.Now()
	v := s.Timeout.Compare(tn)
	return v == -1
}
