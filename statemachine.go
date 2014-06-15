package statemachine

import (
	"fmt"
	"strings"
)

type StateMachine struct {
	current    string
	states     map[stateKey]string
	handlers   map[handlerKey]Handler
	startState func()
}

// NewStateMachine constructs a StateMachine from events and handlers.
//
// The events and states are specified as a slice of Event structs
// specified as Events. Each Event is mapped to one or more internal
// states from Event.Src to Event.Dst.
//
// Handlers are added as a map specified as Handlers where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after eftering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
//
// There are also two short form versions for the most commonly used handlers.
// They are simply the name of the event or state:
//
// 1. <NEW_STATE> - called after entering <NEW_STATE>
//
// 2. <EVENT> - called after event named <EVENT>
//
// If both a shorthand version and a full version is specified it is undefined
// which version of the callback will end up in the internal map. This is due
// to the psuedo random nature of Go maps. No checking for multiple keys is
// currently performed.
func NewStateMachine(initial string, events Events, handlers Handlers) *StateMachine {
	var machine StateMachine
	machine.current = initial
	machine.states = make(map[stateKey]string)
	machine.handlers = make(map[handlerKey]Handler)

	// Build startState map and store sets of all events and states.
	allEvents := make(map[string]bool)
	allStates := make(map[string]bool)
	for _, event := range events {
		for _, src := range event.Src {
			machine.states[stateKey{event.Name, src}] = event.Dst
			allStates[src] = true
			allStates[event.Dst] = true
		}
		allEvents[event.Name] = true
	}

	// Map all handlers to events/states.
	for handlerName, handler := range handlers {
		var target string
		var handlerType handlerType

		switch {
		case strings.HasPrefix(handlerName, "before_"):
			target = strings.TrimPrefix(handlerName, "before_")
			if target == "event" {
				target = ""
				handlerType = beforeEvent
			} else if _, ok := allEvents[target]; ok {
				handlerType = beforeEvent
			}
		case strings.HasPrefix(handlerName, "leave_"):
			target = strings.TrimPrefix(handlerName, "leave_")
			if target == "state" {
				target = ""
				handlerType = leaveState
			} else if _, ok := allStates[target]; ok {
				handlerType = leaveState
			}
		case strings.HasPrefix(handlerName, "enter_"):
			target = strings.TrimPrefix(handlerName, "enter_")
			if target == "state" {
				target = ""
				handlerType = enterState
			} else if _, ok := allStates[target]; ok {
				handlerType = enterState
			}
		case strings.HasPrefix(handlerName, "after_"):
			target = strings.TrimPrefix(handlerName, "after_")
			if target == "event" {
				target = ""
				handlerType = afterEvent
			} else if _, ok := allEvents[target]; ok {
				handlerType = afterEvent
			}
		default:
			target = handlerName
			if _, ok := allStates[target]; ok {
				handlerType = enterState
			} else if _, ok := allEvents[target]; ok {
				handlerType = afterEvent
			}
		}

		if handlerType != noHandler {
			machine.handlers[handlerKey{target, handlerType}] = handler
		}
	}

	return &machine
}

// Current returns the current state of the FSM.
func (machine *StateMachine) Current() string {
	return machine.current
}

// Is returns true if state is the current state.
func (machine *StateMachine) Is(state string) bool {
	return state == machine.current
}

// Can returns true if event can occur in the current state.
func (machine *StateMachine) Can(event string) bool {
	_, ok := machine.states[stateKey{event, machine.current}]
	return ok && (machine.startState == nil)
}

// Can returns true if event can not occure in the current state.
// It is a convenience method to help code read nicely.
func (machine *StateMachine) Cannot(event string) bool {
	return !machine.Can(event)
}

// Event initiates a state startState with the named event.
//
// The call takes a variable number of arguments that will be passed to the
// callback, if defined.
//
// It will return nil if the state change is ok or one of these errors:
//
// - event X inappropriate because previous startState did not complete
//
// - event X inappropriate in current state Y
//
// - event X does not exist
//
// - internal error on state startState
//
// The last error should never occur in this situation and is a sign of an
// internal bug.
func (machine *StateMachine) Event(eventName string, args ...interface{}) error {
	if machine.startState != nil {
		return fmt.Errorf("event %s inappropriate because previous startState did not complete", eventName)
	}

	dst, ok := machine.states[stateKey{eventName, machine.current}]
	if !ok {
		found := false
		for state, _ := range machine.states {
			if state.event == eventName {
				found = true
				break
			}
		}
		if found {
			return fmt.Errorf("event %s inappropriate in current state %s", eventName, machine.current)
		} else {
			return fmt.Errorf("event %s does not exist", eventName)
		}
	}

	if machine.current == dst {
		return nil
	}

	event := &Event{machine, eventName, machine.current, dst, nil, args, false, false}

	// Call the before_ handlers, first the named then the general version.
	if handler, ok := machine.handlers[handlerKey{eventName, beforeEvent}]; ok {
		handler(event)
		if event.canceled {
			return event.Err
		}
	}
	if handler, ok := machine.handlers[handlerKey{"", beforeEvent}]; ok {
		handler(event)
		if event.canceled {
			return event.Err
		}
	}

	machine.startState = func() {
		// Do the state startState.
		machine.current = dst

		// Call the enter_ handlers, first the named then the general version.
		if handler, ok := machine.handlers[handlerKey{machine.current, enterState}]; ok {
			handler(event)
		}
		if handler, ok := machine.handlers[handlerKey{"", enterState}]; ok {
			handler(event)
		}

		// Call the after_ handlers, first the named then the general version.
		if handler, ok := machine.handlers[handlerKey{eventName, afterEvent}]; ok {
			handler(event)
		}
		if handler, ok := machine.handlers[handlerKey{"", afterEvent}]; ok {
			handler(event)
		}
	}

	// Call the leave_ handlers, first the named then the general version.
	if handler, ok := machine.handlers[handlerKey{machine.current, leaveState}]; ok {
		handler(event)
		if event.canceled {
			machine.startState = nil
			return event.Err
		} else if event.async {
			return event.Err
		}
	}
	if handler, ok := machine.handlers[handlerKey{"", leaveState}]; ok {
		handler(event)
		if event.canceled {
			machine.startState = nil
			return event.Err
		} else if event.async {
			return event.Err
		}
	}

	// Perform the rest of the startState, if not asynchronous.
	err := machine.Excute()
	if err != nil {
		return fmt.Errorf("internal error on state startState")
	}

	return event.Err
}

// Excute completes an asynchrounous state change.
//
// The callback for leave_<STATE> must prviously have called Async on its
// event to have initiated an asynchronous state startState.
func (f *StateMachine) Excute() error {
	if f.startState == nil {
		return fmt.Errorf("startState inappropriate because no state change in progress")
	}
	f.startState()
	f.startState = nil
	return nil
}
