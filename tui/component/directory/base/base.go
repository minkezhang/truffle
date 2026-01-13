package base

import (
	"github.com/charmbracelet/bubbletea"
)

type Identifiable interface {
	ID() string
	ParentID() string
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
//	    func() tea.Msg { return m.RegisterMessage{ Node: m } },
//	    tea.Batch( m.child.Init(), ... ),
//	  )
//	}
type RegisterMessage struct {
	Node Identifiable
}

type D struct {
	Nodes    map[string]Identifiable
	Children map[string][]string
	Parent   map[string]string
}

func New() *D {
	return &D{
		Nodes:    map[string]Identifiable{},
		Children: map[string][]string{},
		Parent:   map[string]string{},
	}
}

func (d *D) Init() tea.Cmd { return nil }

func (d *D) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RegisterMessage:
		d.Nodes[msg.Node.ID()] = msg.Node
		d.Children[msg.Node.ParentID()] = append(d.Children[msg.Node.ParentID()], msg.Node.ID())
		d.Parent[msg.Node.ID()] = msg.Node.ParentID()
	}
	return d, nil
}

func (d *D) View() string { return "" }
