package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

type Node struct {
	*focusable.Node
	column   *column.C
	viewport viewport.Model
	node     tea.Model

	clickable_up   *clickable.Node
	clickable_bar  *clickable.Node
	clickable_down *clickable.Node
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
		viewport: viewport.New(o.Column.Content()-1, 0),
		node:     o.Node,
	}
	n.viewport.MouseWheelEnabled = true
	n.clickable_up = clickable.New(n.ID())
	n.clickable_bar = clickable.New(n.ID())
	n.clickable_down = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.node.Init(),
		n.clickable_up.Init(),
		n.clickable_bar.Init(),
		n.clickable_down.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

// scroll_position returns the starting position of the scrollbar of height h.
func (n *Node) scroll_position(h int) int {
	return int(float64(n.viewport.Height-2-h) * n.viewport.ScrollPercent())
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		n.viewport.Height = msg.Height
	case clickable.Click:
		if msg.ID == n.clickable_up.ID() {
			n.viewport, c = n.viewport.Update(tea.KeyMsg{Type: tea.KeyUp})
			cmds = append(cmds, c)
		}
		if msg.ID == n.clickable_bar.ID() {
		}
		if msg.ID == n.clickable_down.ID() {
			n.viewport, c = n.viewport.Update(tea.KeyMsg{Type: tea.KeyDown})
			cmds = append(cmds, c)
		}
	}

	if n.FocusState() == types.FocusStateActive {
		n.viewport, c = n.viewport.Update(msg)
		cmds = append(cmds, c)
	} else {
		switch msg := msg.(type) {
		case clickable.Click:
			if map[string]bool{
				n.clickable_up.ID():   true,
				n.clickable_bar.ID():  true,
				n.clickable_down.ID(): true,
			}[msg.ID] {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
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

	for _, n := range []tea.Model{
		n.clickable_up,
		n.clickable_bar,
		n.clickable_down,
	} {
		_, c := n.Update(msg)
		cmds = append(cmds, c)
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	bar := []string{zone.Mark(n.clickable_up.ID(), "↑")}
	h := 1
	if n.viewport.Height > 0 {
		h = 10 * n.viewport.TotalLineCount() / n.viewport.Height
	}

	scroll_position := n.scroll_position(h)
	for i := 0; i < n.viewport.Height-2; i++ {
		bar = append(bar, map[bool]string{
			false: "░",
			true:  "█",
		}[i == scroll_position || (i > scroll_position && i < scroll_position+h)])
	}
	bar = append(bar, zone.Mark(n.clickable_down.ID(), "↓"))
	scrollbar := lipgloss.JoinVertical(
		lipgloss.Left, bar...,
	)
	if n.viewport.TotalLineCount() <= n.viewport.Height {
		scrollbar = ""
	}
	return n.column.RenderOrDie(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			n.column.WithWidth(n.column.Width()-1).Style().Render(
				n.viewport.View(),
			),
			lipgloss.NewStyle().Foreground(
				color_profile.UIForeground[n.FocusState()],
			).Render(scrollbar),
		),
	)
}
