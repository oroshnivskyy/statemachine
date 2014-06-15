package statemachine

import (
	"fmt"
)

func ExampleNewStateMachine() {
	fsm := NewStateMachine(
		"green",
		Events{
			{Name: "warn", Src: []string{"green"}, Dst: "yellow"},
			{Name: "panic", Src: []string{"yellow"}, Dst: "red"},
			{Name: "panic", Src: []string{"green"}, Dst: "red"},
			{Name: "calm", Src: []string{"red"}, Dst: "yellow"},
			{Name: "clear", Src: []string{"yellow"}, Dst: "green"},
		},
		Handlers{
			"before_warn": func(e *Event) {
				fmt.Println("before_warn")
			},
			"before_event": func(e *Event) {
				fmt.Println("before_event")
			},
			"leave_green": func(e *Event) {
				fmt.Println("leave_green")
			},
			"leave_state": func(e *Event) {
				fmt.Println("leave_state")
			},
			"enter_yellow": func(e *Event) {
				fmt.Println("enter_yellow")
			},
			"enter_state": func(e *Event) {
				fmt.Println("enter_state")
			},
			"after_warn": func(e *Event) {
				fmt.Println("after_warn")
			},
			"after_event": func(e *Event) {
				fmt.Println("after_event")
			},
		},
	)
	fmt.Println(fsm.Current())
	err := fsm.Event("warn")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
}

func ExampleStateMachine_Current() {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	fmt.Println(fsm.Current())
}

func ExampleStateMachine_Is() {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	fmt.Println(fsm.Is("closed"))
	fmt.Println(fsm.Is("open"))
}

func ExampleStateMachine_Can() {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	fmt.Println(fsm.Can("open"))
	fmt.Println(fsm.Can("close"))
}

func ExampleStateMachine_Cannot() {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	fmt.Println(fsm.Cannot("open"))
	fmt.Println(fsm.Cannot("close"))
}

func ExampleStateMachine_Event() {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{},
	)
	fmt.Println(fsm.Current())
	err := fsm.Event("open")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	err = fsm.Event("close")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
}

func ExampleStateMachine_Excute() {
	fsm := NewStateMachine(
		"closed",
		Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		Handlers{
			"leave_closed": func(e *Event) {
				e.Async()
			},
		},
	)
	err := fsm.Event("open")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	err = fsm.Excute()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
}
