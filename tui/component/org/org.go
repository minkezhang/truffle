// Package org tracks application-level stats-related components.
package org

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/focusable"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/log"
)

type Node struct {
	column   *column.C
	children []tea.Model
}

func New(c *column.C) *Node {
	return &Node{
		column: c,
		children: []tea.Model{
			focusable.New(), // directory
			log.New("", column.New(c.Content()).WithBorder(lipgloss.NormalBorder(), true, false, true, false)), // logger
			&errors.Node{}, // error handler
		},
	}
}

func (n *Node) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range n.children {
		cmds = append(cmds, c.Init())
	}
	return tea.Sequence(cmds...)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd
	for i := range n.children {
		n.children[i], c = n.children[i].Update(msg)
		cmds = append(cmds, c)
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	var parts []string
	for _, c := range n.children {
		if v := c.View(); v != "" {
			parts = append(parts, c.View())
		}
	}
	return n.column.RenderOrDie(n.column.Style().Render(lipgloss.JoinVertical(lipgloss.Left, parts...)))
}
