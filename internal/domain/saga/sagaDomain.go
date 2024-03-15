package saga

import (
	"time"
)

type SagaStatus string

type ResponseStatus string

const (
	SUCCESS ResponseStatus = "SUCESS"
	FAILED  ResponseStatus = "FAILED"
)

type Response struct {
	SagaName string         `json:saga_name`
	SagaUUID string         `json:saga_uuid`
	Payload  string         `json:payload`
	Status   ResponseStatus `json:status`
}

const (
	EXECUTING          SagaStatus = "EXECUTING"
	ROLLINGBACK        SagaStatus = "ROLLINGBACK"
	COMPLETED          SagaStatus = "COMPLETED"
	COMPLETED_ROLLBACK SagaStatus = "COMPLETED_ROLLBACK"
)

type Saga struct {
	Id             int64
	SagaName       string
	SagaUUID       string
	Idx            int
	Payload        string
	Status         SagaStatus
	CurrentCommand Command
	LastUpdate     time.Time
}

func NewSaga() *Saga {
	return &Saga{Idx: 0}
}

func (s *Saga) PrepareSaga(response Response) {
	idx, isRollback := handleResponse(response, s.Idx)

	command := GetNextCommanfWithIndex(response.SagaName, s.Idx, isRollback)

	s.CurrentCommand = command
	s.Status = getSagaStatus(isRollback, command)
	s.Payload = response.Payload
	s.Idx = idx
	s.LastUpdate = time.Now()
}

func handleResponse(response Response, idx int) (int, bool) {
	var isRollback bool
	if response.Status == SUCCESS {
		idx++
		isRollback = false
	} else {
		idx--
		isRollback = true
	}
	return idx, isRollback
}

func getSagaStatus(isRollback bool, command Command) SagaStatus {
	if command == "" && isRollback {
		return COMPLETED_ROLLBACK
	} else if command == "" {
		return COMPLETED
	}

	if isRollback {
		return ROLLINGBACK
	} else {
		return EXECUTING
	}
}
