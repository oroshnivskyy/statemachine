package statemachine

type Handlers map[string]Handler

type Handler func(*Event)

type handlerType int

const (
	noHandler handlerType = iota
	beforeEvent
	leaveState
	enterState
	afterEvent
)

// handlerKey is a struct key used for keeping the handlers mapped to a target.
type handlerKey struct {
	// target is either the name of a state or an event depending on which
	target string

	// handlerType is the situation when the callback will be run.
	handlerType handlerType
}
