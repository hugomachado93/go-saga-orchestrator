package statemachine

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

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

func validateDefinition(context *Statemachine) error {
	for _, v := range context.Workflow {
		state := v.State
		event := v.Event
		nextState := v.NextState
		hasState := false
		count := 0
		for _, v := range context.Workflow {
			if v.State == state && v.Event == event {
				count++
			}
			if v.State == nextState {
				hasState = true
				continue
			}
		}

		if count > 1 {
			return fmt.Errorf("can only have one state but more than one were found: %s", state)
		}

		if !v.End && !hasState {
			return fmt.Errorf("the definition is wrong. there is no state for NextState %s", nextState)
		}

		if v.End && v.NextState != "" {
			return fmt.Errorf("final state can`t have nextState")
		}
	}
	return nil
}

func (sm *Statemachine) DrawGraph() ([]byte, error) {
	if err := validateDefinition(sm); err != nil {
		return nil, err
	}

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := graph.Close(); err != nil {
			fmt.Println(err)
		}
		g.Close()
	}()

	for _, v := range sm.Workflow {
		nState := v.NextState
		if v.End {
			nState = "END"
		}

		n1, err := graph.CreateNode(v.State)
		if err != nil {
			return nil, err
		}
		n1.SetFillColor("green")
		n1.SetStyle(cgraph.FilledNodeStyle)
		n2, err := graph.CreateNode(nState)
		if err != nil {
			return nil, err
		}

		if v.End {
			n2.SetFillColor("blue")
		} else {
			n2.SetFillColor("green")
		}
		n2.SetStyle(cgraph.FilledNodeStyle)
		e, err := graph.CreateEdge(v.Event, n1, n2)
		if err != nil {
			return nil, err
		}
		e.SetLabel(v.Event)
	}

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.SVG, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
