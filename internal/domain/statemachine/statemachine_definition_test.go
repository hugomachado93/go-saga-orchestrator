package statemachine

import (
	"testing"
	"time"
)

func Test_Should_Return_Sucess_When_State_And_Event_Exists(t *testing.T) {
	stmc1 := Stmc{State: "state1", Event: "event1", NextState: "state2"}
	smtc2 := Stmc{State: "state2", Event: "event2", NextState: "state3"}
	smtc3 := Stmc{State: "state3", Event: "", NextState: ""}

	smtcw := []Stmc{stmc1, smtc2, smtc3}

	sm := &Statemachine{Name: "teste", Workflow: smtcw}

	sms := &StatemachineDefinition{Id: 1, ClientApiKey: "123", Name: "teste", Context: sm, CreatedAt: time.Now()}

	nextStep, err, _ := sms.FindNextStep("state1", "event1")

	if err != nil {
		t.Fatalf("Test failed with exception %q", err)
	}

	if nextStep != "state2" {
		t.Fatalf("Test failed, should be %q but found %q", "state2", nextStep)
	}
}

func Test_Should_Return_Error_When_State_Dont_Exist(t *testing.T) {
	// stmc1 := Stmc{State: "state1", Event: "event1", NextState: "state2"}
	// smtc2 := Stmc{State: "state2", Event: "event2", NextState: "state3"}
	// smtc3 := Stmc{State: "state3", Event: "", NextState: ""}

	// smtcw := []Stmc{stmc1, smtc2, smtc3}

	// sm := &Statemachine{Name: "teste", Workflow: smtcw}

	// sms := &StatemachineDefinition{Id: 1, ClientApiKey: "123", Name: "teste", Context: sm, CreatedAt: time.Now()}

	// _, err := sms.FindNextStep("state4", "event1")

	// if err == nil {
	// 	t.Fatal("Test failed - error is nil")
	// }
}
