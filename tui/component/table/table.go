package table

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

type Node struct {
	*focusable.Node

	table     table.Model
	clickable *clickable.Node
	data      []string // TODO
}

type O struct {
	Prefix   string
	ParentID string
}

func New(o O) *Node {
	n := &Node{
		Node: focusable.New(o.Prefix, o.ParentID, 1),
		table: table.New(
			table.WithColumns([]table.Column{
				table.Column{
					Title: "Title",
					Width: 50,
				},
				table.Column{
					Title: "Type",
					Width: 20,
				},
				table.Column{
					Title: "Source",
					Width: 10,
				},
				table.Column{
					Title: "Score",
					Width: 10,
				},
			}),
			table.WithFocused(false),
			table.WithHeight(20),
		),
	}
	n.table.SetRows([]table.Row{
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"Frieren", "Anime", "MAL", "★★★☆☆"},
		{"AAAAAA Frieren", "Anime", "MAL", "★★★☆☆"},
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
		lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(
			color_profile.UIForeground[n.FocusState()],
		).Render(
			n.table.View(),
		),
	)
}

func (n *Node) OnFocus(i int) tea.Cmd {
	n.table.Focus()
	return n.Node.OnFocus(i)
}

func (n *Node) OnBlur() tea.Cmd {
	n.table.Blur()
	return n.Node.OnBlur()
}
