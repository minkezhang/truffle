package component_node

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/node/source/view"
	"github.com/minkezhang/truffle/tui/util/node"
)

type O struct {
	ParentID       string
	Column         *column.C
	Node           util_node.N
	CacheDirectory string
}

type Node struct {
	*focusable.Node

	column *column.C
	node   util_node.N
	source *view.Node
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("node-view", o.ParentID, 0),
		column: o.Column,
		node:   o.Node,
	}
	n.source = view.New(view.O{
		ParentID:       n.ID(),
		Column:         o.Column,
		CacheDirectory: o.CacheDirectory,
	})
	return n
}

func (n *Node) SetValue(v util_node.N) tea.Cmd {
	if v == nil {
		return nil
	}
	n.node = v
	source, err := v.Virtual()
	if err != nil {
		return func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelWarn,
				fmt.Sprintf("%v: cannot get a merged source component: %v", err),
			)
		}
	}
	return n.source.SetValue(source)
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.source.Init(),
		n.SetValue(n.node),
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

	_, c = n.source.Update(msg)
	cmds = append(cmds, c)

	switch msg := msg.(type) {
	case message.GetNodeResponseMessage:
		cmds = append(
			cmds,
			n.SetValue(msg.Body.Value),
			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: received GetNodeResponseMessage: %v", n.ID(), msg),
				)
			},
		)
	case message.PutResponseMessage:
		cmds = append(
			cmds,
			n.SetValue(msg.Body.Value.Node),
			func() tea.Msg {
				return errors.ToLogMessage(
					errors.LevelDebug,
					fmt.Sprintf("%v: received PutResponseMessage: %v", n.ID(), msg),
				)
			},
		)
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	if n.node == nil {
		return ""
	}
	return n.column.RenderOrDie(n.column.Style().Render(
		n.source.View(),
	))
}
