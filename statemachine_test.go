package statemachine

import (
	"fmt"
	"testing"
)

func TestSameState(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "start"},
		},
		Handlers{},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestInappropriateEvent(t *testing.T) {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	err := fsm.Event("close")
	if err.Error() != "event close inappropriate in current state closed" {
		t.FailNow()
	}
}

func TestInvalidEvent(t *testing.T) {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	err := fsm.Event("lock")
	if err.Error() != "event lock does not exist" {
		t.FailNow()
	}
}

func TestMultipleSources(t *testing.T) {
	fsm := NewStateMachine(
		"one",
		Events{
			{Name: "first", Src: []string{"one"}, Dst: "two"},
			{Name: "second", Src: []string{"two"}, Dst: "three"},
			{Name: "reset", Src: []string{"one", "two", "three"}, Dst: "one"},
		},
		Handlers{},
	)

	fsm.Event("first")
	if fsm.Current() != "two" {
		t.FailNow()
	}
	fsm.Event("reset")
	if fsm.Current() != "one" {
		t.FailNow()
	}
	fsm.Event("first")
	fsm.Event("second")
	if fsm.Current() != "three" {
		t.FailNow()
	}
	fsm.Event("reset")
	if fsm.Current() != "one" {
		t.FailNow()
	}
}

func TestMultipleEvents(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "first", Src: []string{"start"}, Dst: "one"},
			{Name: "second", Src: []string{"start"}, Dst: "two"},
			{Name: "reset", Src: []string{"one"}, Dst: "reset_one"},
			{Name: "reset", Src: []string{"two"}, Dst: "reset_two"},
			{Name: "reset", Src: []string{"reset_one", "reset_two"}, Dst: "start"},
		},
		Handlers{},
	)

	fsm.Event("first")
	fsm.Event("reset")
	if fsm.Current() != "reset_one" {
		t.FailNow()
	}
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.FailNow()
	}

	fsm.Event("second")
	fsm.Event("reset")
	if fsm.Current() != "reset_two" {
		t.FailNow()
	}
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestGenericHandlers(t *testing.T) {
	beforeEvent := false
	leaveState := false
	enterState := false
	afterEvent := false

	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"before_event": func(e *Event) {
				beforeEvent = true
			},
			"leave_state": func(e *Event) {
				leaveState = true
			},
			"enter_state": func(e *Event) {
				enterState = true
			},
			"after_event": func(e *Event) {
				afterEvent = true
			},
		},
	)

	fsm.Event("run")
	if !(beforeEvent && leaveState && enterState && afterEvent) {
		t.FailNow()
	}
}

func TestSpecificHandlers(t *testing.T) {
	beforeEvent := false
	leaveState := false
	enterState := false
	afterEvent := false

	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"before_run": func(e *Event) {
				beforeEvent = true
			},
			"leave_start": func(e *Event) {
				leaveState = true
			},
			"enter_end": func(e *Event) {
				enterState = true
			},
			"after_run": func(e *Event) {
				afterEvent = true
			},
		},
	)

	fsm.Event("run")
	if !(beforeEvent && leaveState && enterState && afterEvent) {
		t.FailNow()
	}
}

func TestSpecificHandlersShortform(t *testing.T) {
	enterState := false
	afterEvent := false

	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"end": func(e *Event) {
				enterState = true
			},
			"run": func(e *Event) {
				afterEvent = true
			},
		},
	)

	fsm.Event("run")
	if !(enterState && afterEvent) {
		t.FailNow()
	}
}

func TestCancelBeforeGenericEvent(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"before_event": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestCancelBeforeSpecificEvent(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"before_run": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestCancelLeaveGenericState(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"leave_state": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestCancelLeaveSpecificState(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"leave_start": func(e *Event) {
				e.Cancel()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestAsyncExcuteGenericState(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"leave_state": func(e *Event) {
				e.Async()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
	fsm.Excute()
	if fsm.Current() != "end" {
		t.FailNow()
	}
}

func TestAsyncExcuteSpecificState(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"leave_start": func(e *Event) {
				e.Async()
			},
		},
	)
	fsm.Event("run")
	if fsm.Current() != "start" {
		t.FailNow()
	}
	fsm.Excute()
	if fsm.Current() != "end" {
		t.FailNow()
	}
}

func TestAsyncExcuteInProgress(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
			{Name: "reset", Src: []string{"end"}, Dst: "start"},
		},
		Handlers{
			"leave_start": func(e *Event) {
				e.Async()
			},
		},
	)
	fsm.Event("run")
	err := fsm.Event("reset")
	if err.Error() != "event reset inappropriate because previous startState did not complete" {
		t.FailNow()
	}
	fsm.Excute()
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.FailNow()
	}
}

func TestAsyncExcuteNotInProgress(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
			{Name: "reset", Src: []string{"end"}, Dst: "start"},
		},
		Handlers{},
	)
	err := fsm.Excute()
	if err.Error() != "startState inappropriate because no state change in progress" {
		t.FailNow()
	}
}

func TestHandlerNoError(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"run": func(e *Event) {
			},
		},
	)
	e := fsm.Event("run")
	if e != nil {
		t.FailNow()
	}
}

func TestHandlerError(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"run": func(e *Event) {
				e.Err = fmt.Errorf("error")
			},
		},
	)
	e := fsm.Event("run")
	if e.Error() != "error" {
		t.FailNow()
	}
}

func TestHandlerArgs(t *testing.T) {
	fsm := NewStateMachine(
		"start",
		Events{
			{Name: "run", Src: []string{"start"}, Dst: "end"},
		},
		Handlers{
			"run": func(e *Event) {
				if len(e.Args) != 1 {
					t.Fatal("too few arguments")
				}
				arg, ok := e.Args[0].(string)
				if !ok {
					t.Fatal("not a string argument")
				}
				if arg != "test" {
					t.Fatal("incorrect argument")
				}
			},
		},
	)
	fsm.Event("run", "test")
}
