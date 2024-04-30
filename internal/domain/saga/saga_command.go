package saga

import "time"

type SagaCommand struct {
	Id             int64      `db:"id"`
	Payload        string     `db:"payload"`
	CurrentState   string     `db:"state"`
	Status         Status     `db:"status"`
	Timeout        *time.Time `db:"timeout"`
	CreatedAt      time.Time  `db:"created_at"`
	DelayedMessage *time.Time `db:"delayed_message"`
	LastUpdate     time.Time  `db:"last_update"`
	SagaId         *int64     `db:"saga_id"`
}

type Status string

const (
	WAITING_PROCESS   Status = "WAITING_PROCESS"
	REQUEST_SENT      Status = "REQUEST_SENT"
	RESPONSE_RECEIVED Status = "RESPONSE_RECEIVED"
)

func NewCommand(state string, sagaId *int64, delay *int) SagaCommand {
	var dm time.Time
	if delay != nil {
		dm = time.Now().Add(time.Second * time.Duration(*delay))
	}
	return SagaCommand{CurrentState: state, Status: WAITING_PROCESS, CreatedAt: time.Now(), LastUpdate: time.Now(), SagaId: sagaId, DelayedMessage: &dm}
}

func (sc *SagaCommand) RequestCommand() {
	t := time.Now().Add(time.Duration(time.Duration.Seconds(120)))

	sc.Status = REQUEST_SENT
	sc.LastUpdate = time.Now()
	sc.Timeout = &t
}

func (sc *SagaCommand) ReceiveResponse() {
	if sc.Status == REQUEST_SENT {
		sc.Status = RESPONSE_RECEIVED
	}
	sc.Status = REQUEST_SENT
}

func (sc *SagaCommand) HasCommandTimeout() bool {
	tn := time.Now()
	v := sc.Timeout.Compare(tn)
	return v == -1
}
