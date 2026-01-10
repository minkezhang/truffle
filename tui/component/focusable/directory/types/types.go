package types

import (
	"github.com/charmbracelet/bubbletea"
)

type FocusState int

const (
	FocusStateNone FocusState = iota
	FocusStateActive
)

type Node interface {
	tea.Model

	Focusable
	Identifiable
	Renderable
}

type Renderable interface {
}

type Identifiable interface {
	ID() string
	ParentID() string
}

type Focusable interface {
	// OnFocus may be called in the Update method by the directory. This
	// function directly instructs the node to activate any sub-elements for
	// input, e.g. focusing a specific input field.
	OnFocus(i int) tea.Cmd

	// OnBlur may be called in the Update method by the directory. This
	// function directly instructs the node to deactivate all sub-elements.
	OnBlur() tea.Cmd

	// NElements returns the number of interactable (and therefore
	// tab-focusable) elements in this node. NElements does not take into
	// account the number of [Node] children.
	NElements() int

	FocusIndex() int
	FocusState() FocusState
}

// RegisterNodeMessage should be fired by all [Node] implementations during
// Init. This message will be handled by
// [github.com/minkezhang/truffle/tui/components/directory.D] which will add the
// node to a central repo.
//
// Init should be of the form
//
//	func (m M) Init() tea.Cmd {
//	  return tea.Sequence(
//	    func() tea.Msg { return m.RegisterNodeMessage{ Node: m } },
//	    tea.Batch( m.child.Init(), ... ),
//	  )
//	}
type RegisterNodeMessage struct {
	Node Node
}

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
