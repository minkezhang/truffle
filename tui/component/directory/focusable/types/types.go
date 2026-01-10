package types

type FocusState int

const (
	FocusStateNone FocusState = iota
	FocusStateActive
)

// EOFMessage will be sent by a [Focusable] instance in the OnFocus function.
// This will be handled by
// [github.com/minkezhang/truffle/tui/components/directory.D] and pass focus
// onto the next (or previous) node.
type EOFMessage struct {
	ID string
	// IsHead indicates if focus should be handed off to the previous or
	// next node in the focus traversal order.
	IsHead bool
}

// FocusMessage will be sent by [Node] implementations to explicitly take
// control of user input. The directory will hand control as part of the message
// handling by calling OnFocus.
type FocusMessage struct {
	ID    string
	Index int
}
