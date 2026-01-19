// Package display encapsulates all individual display components.
//
// This component will be encapsulated in a viewport.
package display

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/db/message"
	"github.com/minkezhang/truffle/tui/component/org"
	"github.com/minkezhang/truffle/tui/component/search"
	"github.com/minkezhang/truffle/tui/component/table"
	"github.com/minkezhang/truffle/tui/util/form"
)

type O struct {
	Column         *column.C
	CacheDirectory string
}

func Make(o O) Node {
	return Node{
		org: org.New(o.Column),
		search_bar: search.New(search.O{
			Column: o.Column,
		}),
		table: table.New(table.O{
			Prefix: "table-view",
			Column: o.Column,
			Key: form.Key{
				Label: "",
				Key:   "table-selection",
			},
			CacheDirectory: o.CacheDirectory,
		}),
	}
}

type Node struct {
	org        *org.Node
	search_bar *search.Node
	table      *table.Node
}

func (n Node) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range []tea.Model{
		n.org,
		n.search_bar,
		n.table,
	} {
		cmds = append(cmds, c.Init())
	}
	return tea.Sequence(cmds...)
}

func (n Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Reroute messages.
	switch msg := msg.(type) {
	case message.SearchResponseMessage:
		if msg.ID == n.search_bar.ID() {
			cmds = append(cmds, func() tea.Msg {
				return message.SearchResponseMessage{
					ID:      n.table.ID(),
					Results: msg.Results,
				}
			})
		}
	}

	for _, c := range []tea.Model{
		n.org,
		n.search_bar,
		n.table,
	} {
		_, d := c.Update(msg)
		cmds = append(cmds, d)
	}
	return n, tea.Batch(cmds...)
}

func (n Node) View() string {
	parts := []string{
		n.org.View(),
		n.search_bar.View(),
	}
	if !n.table.IsInvisible() {
		parts = append(parts, n.table.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
