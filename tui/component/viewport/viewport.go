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
		viewport: viewport.New(o.Column.Content()-2, 0),
		node:     o.Node,
	}
	n.viewport.MouseWheelEnabled = true
	n.clickable_up = clickable.New(n.ID())
	n.clickable_down = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.node.Init(),
		n.clickable_up.Init(),
		n.clickable_down.Init(),
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
	case clickable.Click:
		if msg.ID == n.clickable_up.ID() {
			n.viewport, c = n.viewport.Update(tea.KeyMsg{Type: tea.KeyUp})
			cmds = append(cmds, c)
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
			if msg.Type == tea.KeyEsc && n.viewport.Height < n.viewport.TotalLineCount() {
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
		n.clickable_down,
	} {
		_, c := n.Update(msg)
		cmds = append(cmds, c)
	}

	return n, tea.Sequence(cmds...)
}

// scroll_position returns the starting position of the scrollbar of height h.
func (n *Node) scroll_position(h int) int {
	return int(float64(n.viewport.Height-2-h) * n.viewport.ScrollPercent())
}

func (n *Node) scroll_height() int {
	h := 1
	if n.viewport.TotalLineCount() > 0 {
		h = int(float64(n.viewport.Height-2) * float64(n.viewport.Height) / float64(n.viewport.TotalLineCount()))
		if h == 0 {
			h = 1
		}
	}
	return h
}

func (n *Node) IsInvisible() bool {
	return n.Node.IsInvisible() || n.viewport.Height >= n.viewport.TotalLineCount()
}

func (n *Node) View() string {
	scroll_height := n.scroll_height()
	scroll_position := n.scroll_position(scroll_height)

	bar := []string{}
	for i := 0; i < n.viewport.Height-2; i++ {
		bar = append(bar, map[bool]string{
			false: lipgloss.NewStyle().Foreground(
				color_profile.ForegroundNegligible,
			).Render("░"),
			true: lipgloss.NewStyle().Foreground(
				color_profile.UIForeground[n.FocusState()],
			).Render("█"),
		}[i >= scroll_position && i < scroll_position+scroll_height])
	}

	scrollbar := []string{
		zone.Mark(n.clickable_up.ID(), "↑"),
		lipgloss.JoinVertical(
			lipgloss.Left,
			bar...,
		),
		zone.Mark(n.clickable_down.ID(), "↓"),
	}
	if n.IsInvisible() {
		scrollbar = []string{}
	}
	return n.column.RenderOrDie(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			n.column.WithWidth(n.column.Width()-2).Style().Render(
				n.viewport.View(),
			),
			lipgloss.NewStyle().Margin(0, 0, 0, 1).Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					scrollbar...,
				),
			),
		),
	)
}
