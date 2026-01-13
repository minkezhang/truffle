package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

type Node struct {
	*focusable.Node
	column    *column.C
	viewport  viewport.Model
	node      tea.Model
}

type O struct {
	Prefix   string
	ParentID string
	Column   *column.C
	Node     tea.Model
}

func New(o O) *Node {
	n := &Node{
		Node:     focusable.New(o.Prefix, o.ParentID, 1),
		column:   o.Column,
		viewport: viewport.New(o.Column.Content(), 0),
		node:     o.Node,
	}
	n.viewport.MouseWheelEnabled = true
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.node.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		n.viewport.Height = msg.Height
	}

	if n.FocusState() == types.FocusStateActive {
		n.viewport, c = n.viewport.Update(msg)
		cmds = append(cmds, c)
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEsc {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
		}
	}

	n.node, c = n.node.Update(msg)
	cmds = append(cmds, c)

	n.viewport.SetContent(n.node.View())

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string { return n.viewport.View() }
