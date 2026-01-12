package log

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

const (
	bufsize = 100
)

type Node struct {
	*focusable.Node

	viewport  viewport.Model
	clickable *clickable.Node
	column    *column.C
	lines     []string
}

func New(parent_id string, c *column.C) *Node {
	n := &Node{
		Node:     focusable.New("log", parent_id, 1),
		viewport: viewport.New(c.Content(), 15),
		column:   c,
		lines:    []string{},
	}
	n.clickable = clickable.New(n.ID())
	n.viewport.MouseWheelEnabled = true
	n.viewport.SetHorizontalStep(5)
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
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
	case errors.LogMessage:
		parts := strings.Split(msg.V, "\n")
		header := fmt.Sprintf("%v(%v):", strings.ToUpper(msg.L.String()), msg.T.Format("15:04:05"))
		lines := []string{}
		for i, p := range parts {
			var l string
			if i == 0 {
				l = fmt.Sprintf("%v %v", header, p)
			} else {
				l = fmt.Sprintf("%v %v", strings.Repeat(" ", len(header)), p)
			}
			lines = append(
				lines,
				lipgloss.NewStyle().Foreground(color_profile.LogForeground[msg.L]).Render(l),
			)
		}
		n.lines = append(lines, n.lines...)
		if len(n.lines) > bufsize {
			n.lines = n.lines[:bufsize]
		}
		n.viewport.SetContent(strings.Join(n.lines, "\n"))
	}

	if n.FocusState() == types.FocusStateActive {
		n.viewport, c = n.viewport.Update(msg)
		cmds = append(cmds, c)
	} else {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
		}
		_, c = n.clickable.Update(msg)
		cmds = append(cmds, c)
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return n.column.RenderOrDie(
		zone.Mark(
			n.clickable.ID(),
			n.column.Style().BorderForeground(
				color_profile.UIForeground[n.FocusState()],
			).Render(
				lipgloss.JoinVertical(
					lipgloss.Right,
					n.viewport.View(),
					lipgloss.NewStyle().Foreground(color_profile.SupplementaryText).Render(
						fmt.Sprintf("%3.f%%", n.viewport.ScrollPercent()*100),
					),
				),
			),
		),
	)
}
