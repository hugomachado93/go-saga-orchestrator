package requests

type Stmc struct {
	State     string `json:"state"`
	Event     string `json:"event"`
	NextState string `json:"nextState"`
}

type Statemachine struct {
	Name     string `json:"name"`
	Workflow []Stmc `json:"workflow"`
}
