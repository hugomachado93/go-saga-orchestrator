package statemachine

type Stmc struct {
	State     string `json:"state"`
	Event     string `json:"event"`
	NextState string `json:"nextState"`
	End       bool   `json:"end"`
}

type Statemachine struct {
	Name     string `json:"name"`
	Workflow []Stmc `json:"workflow"`
}
