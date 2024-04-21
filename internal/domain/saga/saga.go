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
	Id             int64      `db:"id"`
	ApiKey         string     `db:"api_key"`
	SagaName       string     `db:name`
	SagaUUID       string     `db:"uuid"`
	Payload        string     `db:"payload"`
	CurrentState   string     `db:"state"`
	Status         Status     `db:"status"`
	Timeout        *time.Time `db:"timeout"`
	CreatedAt      time.Time  `db:"created_at"`
	DelayedMessage *time.Time `db:"delayed_message"`
	LastUpdate     time.Time  `db:"last_update"`
}

func NewSaga(apiKey string, uuid string, name string) *Saga {
	return &Saga{ApiKey: apiKey, SagaUUID: uuid, SagaName: name, CreatedAt: time.Now()}
}

func (s *Saga) PrepareNextCommand(payload string, nextState string, delay *int) {
	s.CurrentState = nextState
	s.Payload = payload
	s.Status = WAITING_PROCESS
	s.LastUpdate = time.Now()
	s.CreatedAt = time.Now()

	if delay != nil {
		t := time.Now().Add(time.Second * time.Duration(*delay))
		s.DelayedMessage = &t
	}

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
