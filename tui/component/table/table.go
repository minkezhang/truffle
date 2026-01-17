package table

import (
	"github.com/charmbracelet/bubbles/table"
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

var (
	styles = map[types.FocusState]table.Styles{
		types.FocusStateNone: table.Styles{
			Header: table.DefaultStyles().Header.Foreground(
				color_profile.UIForeground[types.FocusStateNone],
			).Border(
				lipgloss.NormalBorder(), false, false, true, false,
			).BorderForeground(
				color_profile.UIForeground[types.FocusStateNone],
			),
			Cell: lipgloss.NewStyle().Padding(0, 1),
			Selected: lipgloss.NewStyle().Foreground(
				color_profile.ForegroundNormal,
			).Background(
				color_profile.BackgroundNegligible,
			),
		},
		types.FocusStateActive: table.Styles{
			Header: table.DefaultStyles().Header.Foreground(
				color_profile.UIForeground[types.FocusStateActive],
			).Border(
				lipgloss.NormalBorder(), false, false, true, false,
			).BorderForeground(
				color_profile.UIForeground[types.FocusStateActive],
			),
			Cell: lipgloss.NewStyle().Padding(0, 1),
			Selected: lipgloss.NewStyle().Foreground(
				color_profile.ForegroundNegligible,
			).Background(
				color_profile.ForegroundNormal,
			),
		},
	}
)

type Node struct {
	*focusable.Node

	column    *column.C
	table     table.Model
	clickable *clickable.Node
	data      []string // TODO
}

type O struct {
	Prefix   string
	ParentID string
	Column   *column.C
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New(o.Prefix, o.ParentID, 1),
		column: o.Column,
		table: table.New(
			table.WithColumns([]table.Column{
				table.Column{
					Title: "Title",
					Width: o.Column.Content() - /* other columns */ 45 - /* padding */ 10 - /* border-left */ 1,
				},
				table.Column{
					Title: "Media", // e.g. "Light Novel"
					Width: 15,
				},
				table.Column{
					Title: "API",
					Width: 15,
				},
				table.Column{
					Title: "Status",
					Width: 10,
				},
				table.Column{
					Title: "Score",
					Width: 5,
				},
			}),
			table.WithFocused(false),
			table.WithHeight(20),
		),
	}
	n.table.SetStyles(styles[n.FocusState()])
	n.table.SetRows([]table.Row{
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
		{"Frieren", "Light Novel", "MAL", "Queued", "★★★☆☆"},
	})
	n.clickable = clickable.New(n.ID())
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

	_, c = n.clickable.Update(msg)
	cmds = append(cmds, c)

	if n.FocusState() == types.FocusStateActive {
		n.table, c = n.table.Update(msg)
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
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(
		n.clickable.ID(),
		lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true, false, true, false).BorderForeground(
			color_profile.UIForeground[n.FocusState()],
		).Render(
			n.table.View(),
		),
	)
}

func (n *Node) OnFocus(i int) tea.Cmd {
	n.table.SetStyles(styles[types.FocusStateActive])
	n.table.Focus()
	return n.Node.OnFocus(i)
}

func (n *Node) OnBlur() tea.Cmd {
	n.table.SetStyles(styles[types.FocusStateNone])
	n.table.Blur()
	return n.Node.OnBlur()
}
