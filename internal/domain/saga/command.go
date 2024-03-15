package saga

type Command string

const (
	AUTHORIZE_PAYMENT Command = "AUTHORIZE_PAYMENT"
	COMPLETE_PAYMENT  Command = "COMPLETE_PAYMENT"
)
