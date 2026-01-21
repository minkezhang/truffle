package view

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/image"
	"github.com/minkezhang/truffle/tui/util/node"
)

const (
	image_width = 50
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
	image  *image.Node
	node   util_node.N
	source source.S
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("node-view", o.ParentID, 0),
		column: o.Column,
		node:   o.Node,
	}
	n.image = image.New(image.O{
		ParentID:       n.ID(),
		URL:            "",
		Width:          image_width,
		CacheDirectory: o.CacheDirectory,
	})
	return n
}

func (n *Node) SetValue(v util_node.N) tea.Cmd {
	if v == nil {
		return nil
	}
	return tea.Sequence(
		func() tea.Msg {
			source, err := v.Virtual()
			if err != nil {
				n.node = nil
				return errors.ToLogMessage(
					errors.LevelWarn,
					fmt.Sprintf("%v: cannot get source from node: %v", n.ID(), err),
				)
			}

			n.node = v
			n.source = source
			return n.image.SetValue(n.source.PreviewURL())()
		},
	)
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.image.Init(),
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

	_, c = n.image.Update(msg)
	cmds = append(cmds, c)

	switch msg := msg.(type) {
	case message.GetNodeResponseMessage:
		cmds = append(
			cmds,
			tea.Batch(
				func() tea.Msg {
					return errors.ToLogMessage(
						errors.LevelDebug,
						fmt.Sprintf("%v: got node response message: %v", n.ID(), msg),
					)
				},
				n.SetValue(msg.Body.Value),
			),
		)
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	if n.node == nil {
		return ""
	}

	return n.column.RenderOrDie(
		n.column.Style().Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				n.image.View(),
				n.column.WithWidth(n.column.Width()-image_width).Style().Render(
					lipgloss.JoinVertical(
						lipgloss.Left,
						n.source.Title().Title(),
						n.source.Synopsis(),
					),
				),
			),
		),
	)
}
