package list_item

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type Node struct {
	*focusable.Node
	column *column.C

	value     form.Value[node.N]
	clickable *clickable.Node
}

type O struct {
	Value  form.Value[node.N]
	Parent string
	Column *column.C
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("node-list-item", o.Parent, 1),
		value:  o.Value,
		column: o.Column.WithWidth(150).WithBorder(lipgloss.NormalBorder(), true, false, true, false),
	}
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

func (n *Node) Value() form.Value[node.N] { return n.value }

type HighlightListItemMessage struct {
	ID    string
	Value form.Value[node.N]
}

type SubmitListItemMessage HighlightListItemMessage

func (n *Node) submit() tea.Cmd {
	m := SubmitListItemMessage{
		ID:    n.ID(),
		Value: n.Value(),
	}
	return tea.Sequence(
		func() tea.Msg { return m },
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: submitted message %v", n.ID(), m),
			)
		},
	)

}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				cmds = append(cmds, n.submit())
			}
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, n.submit())
			}
		}
	} else {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, tea.Sequence(
					func() tea.Msg {
						return types.FocusMessage{
							ID:    n.ID(),
							Index: 0,
						}
					},
					n.submit(),
				))
			}
		}

	}
	_, c = n.clickable.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func truncate(s string, w int) string {
	return ansi.Truncate(s, w, "…")
}

var (
	columns = []struct {
		header  string
		width   int
		view    func(data node.N) string
		special bool
	}{
		{
			header: "title",
			width:  50,
			view:   func(data node.N) string { return "Frieren" },
		},
		{
			header:  "source-type",
			width:   11,
			view:    func(data node.N) string { return "Light Novel" },
			special: true,
		},
		{
			header: "api",
			width:  12,
			view:   func(data node.N) string { return "Truffle" },
		},
		{
			header: "status",
			width:  9,
			view:   func(data node.N) string { return "Queued" },
		},
		{
			header: "score",
			width:  5,
			view:   func(data node.N) string { return "★★★☆☆" },
		},
	}
)

func (n *Node) View() string {
	style := lipgloss.NewStyle().Width(
		n.column.Content(),
	).Border(
		map[types.FocusState]lipgloss.Border{
			types.FocusStateActive: lipgloss.ThickBorder(),
			types.FocusStateNone:   lipgloss.HiddenBorder(),
		}[n.FocusState()],
		false, false, false, true,
	).BorderForeground(
		color_profile.UIForeground[n.FocusState()],
	).Foreground(
		color_profile.UIForeground[n.FocusState()],
	)

	var parts []string
	for _, c := range columns {
		d := column.New(c.width + 2)
		parts = append(
			parts,
			d.RenderOrDie(
				lipgloss.NewStyle().Padding(0, 1).Background(
					map[types.FocusState]lipgloss.TerminalColor{
						types.FocusStateActive: color_profile.BackgroundNegligible,
						types.FocusStateNone:   style.GetBackground(),
					}[n.FocusState()],
				).Bold(
					n.FocusState() == types.FocusStateActive,
				).Render(
					d.Style().Width(c.width).MaxWidth(c.width).Inline(true).Render(
						ansi.Truncate(c.view(n.Value().Value), c.width, "…"),
					),
				),
			),
		)
	}

	return n.column.RenderOrDie(n.column.Style().Render(
		zone.Mark(
			n.clickable.ID(),
			style.Render(lipgloss.JoinHorizontal(lipgloss.Top, parts...)),
		),
	))
}

func (n *Node) OnFocus(i int) tea.Cmd {
	return tea.Sequence(
		n.Node.OnFocus(i),
		func() tea.Msg {
			return HighlightListItemMessage{
				ID:    n.ID(),
				Value: n.Value(),
			}
		},
	)
}
