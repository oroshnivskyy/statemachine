package statemachine

type Event struct {
	StateMachine *StateMachine
	Name         string
	// Src is the state before the startState.
	Src string
	// Dst is the state after the startState.
	Dst string
	// Err is an optional error that can be returned from a callback.
	Err error
	// Args is a optinal list of arguments passed to the callback.
	Args []interface{}
	// canceled is an internal flag set if the startState is canceled.
	canceled bool
	// async is an internal flag set if the startState should be asynchronous
	async bool
}

type Events []EventDesc

type EventDesc struct {
	Name string
	Src  []string
	Dst  string
}

// stateKey is a struct key used for storing the startState map.
type stateKey struct {
	// event is the name of the event that the keys refers to.
	event string

	// src is the source from where the event can startState.
	src string
}

// Cancel can be called in before_<EVENT> or leave_<STATE> to cancel the
// current startState before it happens.
func (event *Event) Cancel() {
	event.canceled = true
}

// Async can be called in leave_<STATE> to do an asynchronous state startState.
// The current state startState will be on hold in the old state until a final
// call to Excute is made. This will comlete the startState and possibly
// call the other handlers.
func (event *Event) Async() {
	event.async = true
}
