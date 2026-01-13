// Package org tracks application-level stats-related components.
package org

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/log"
	"github.com/minkezhang/truffle/tui/util/color_profile"

	directory "github.com/minkezhang/truffle/tui/component/directory/focusable"
)

type Node struct {
	*focusable.Node

	column      *column.C
	children    []tea.Model
	is_expanded bool
	clickable   *clickable.Node
}

func New(c *column.C) *Node {
	n := &Node{
		Node:   focusable.New("org", "", 1),
		column: c,
		children: []tea.Model{
			directory.New(), // directory
			&errors.Node{},  // error handler
		},
	}
	n.children = append(
		n.children,
		// logger
		log.New(n.ID(), column.New(c.Content()).WithBorder(
			lipgloss.NormalBorder(), false, false, true, false),
		),
	)
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range n.children {
		cmds = append(cmds, c.Init())
	}
	cmds = append(cmds,
		n.clickable.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
	for _, c := range n.children {
		if c, ok := c.(directory.Node); ok {
			cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
		}
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if n.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				n.is_expanded = !n.is_expanded
				for _, c := range n.children {
					if c, ok := c.(directory.Node); ok {
						cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
					}
				}
			}
		}
	case clickable.Click:
		if msg.ID == n.clickable.ID() {
			n.is_expanded = !n.is_expanded
			cmds = append(cmds, func() tea.Msg {
				return types.FocusMessage{
					ID: n.ID(),
				}
			})
			for _, c := range n.children {
				if c, ok := c.(directory.Node); ok {
					cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
				}
			}
		}
	}

	_, c = n.clickable.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	header := lipgloss.NewStyle().Foreground(
		color_profile.UIForeground[n.FocusState()],
	).Render(
		fmt.Sprintf(
			"%s (debug)",
			zone.Mark(
				n.clickable.ID(),
				map[bool]string{
					true:  "↓",
					false: "→",
				}[n.is_expanded],
			),
		),
	)
	parts := []string{header}
	if n.is_expanded {
		for _, c := range n.children {
			if v := c.View(); v != "" {
				parts = append(parts, c.View())
			}
		}
	}
	return n.column.RenderOrDie(
		n.column.Style().Render(
			lipgloss.JoinVertical(lipgloss.Left, parts...),
		),
	)
}
