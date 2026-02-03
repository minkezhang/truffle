// Package display encapsulates all individual display components.
//
// This component will be encapsulated in a viewport.
package display

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/dropdown"
	"github.com/minkezhang/truffle/tui/component/node"
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
		dropdown: dropdown.New(dropdown.O{
			Prefix: "dropdown-test",
			Width:  50,
			Key:    form.Key{"Type", "type"},
			Choices: []form.Key{
				form.Key{"Anime", "anime"},
				form.Key{"0", "book"},
				form.Key{"01", "book"},
				form.Key{"012", "book"},
				form.Key{"0123", "book"},
				form.Key{"01234", "book"},
				form.Key{"012345", "book"},
				form.Key{"0123456", "book"},
				form.Key{"01234567", "book"},
				form.Key{"012345678", "book"},
				form.Key{"0123456789", "book"},
				form.Key{"01234567890", "book"},
			},
			Prompt: "┃ ",
			Height: 5,
		}),
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
		node: component_node.New(component_node.O{
			CacheDirectory: o.CacheDirectory,
			Column:         o.Column,
		}),
	}
}

type Node struct {
	org        *org.Node
	search_bar *search.Node
	table      *table.Node
	node       *component_node.Node
	dropdown   *dropdown.Node
}

func (n Node) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range []tea.Model{
		n.org,
		n.search_bar,
		n.table,
		n.node,
		n.dropdown,
	} {
		cmds = append(cmds, c.Init())
	}
	return tea.Sequence(cmds...)
}

func (n Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	for _, c := range []tea.Model{
		n.org,
		n.search_bar,
		n.table,
		n.node,
		n.dropdown,
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
		parts = append(parts, lipgloss.NewStyle().Margin(1, 0, 0, 0).Render(n.table.View()))
	}

	parts = append(parts, n.node.View())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, parts...),
		n.dropdown.View(),
	)
}
