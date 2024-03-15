package outboxsaga

import "time"

type OutboxState string

const (
	PROCESSING OutboxState = "PROCESSING"
	PROCESSED  OutboxState = "PROCESSED"
)

type Outboxsaga struct {
	payload     string
	outboxState OutboxState
	registered  time.Time
}
