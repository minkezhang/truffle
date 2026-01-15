package search

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/checkbox_group"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/util/form"
)

type Node struct {
	*focusable.Node

	column *column.C
	input  *textinput.Node

	apis         *checkbox_group.Node
	source_types *checkbox_group.Node
	options *checkbox_group.Node
}

type O struct {
	ParentID string
	Column   *column.C
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("search", o.ParentID, 0),
		column: o.Column,
	}
	n.input = textinput.New(textinput.O{
		Prefix:      "search-textinput",
		ParentID:    n.ID(),
		Width:       o.Column.Content(),
		Placeholder: "Frieren, mal:manga/52991",
		Prompt:      "⚲ ",
		Value: form.Value[string]{
			Key: form.Key{
				Label: "Search",
				Key:   "search-textinput",
			},
		},
	})
	n.apis = checkbox_group.New(checkbox_group.O{
		Prefix:   "search-apis",
		ParentID: n.ID(),
		IsRadio:  false,
		Value: form.Value[[]form.Value[bool]]{
			Key: form.Key{
				Label: "APIs",
				Key:   "search-apis",
			},
			Value: []form.Value[bool]{
				form.Value[bool]{
					Key: form.Key{
						Label: "Truffle",
						Key:   "truffle",
					},
					Value: true,
				},
				form.Value[bool]{
					Key: form.Key{
						Label: "MAL",
						Key:   "mal",
					},
					Value: true,
				},
			},
		},
	})
	n.source_types = checkbox_group.New(checkbox_group.O{
		Prefix:   "search-source-types",
		ParentID: n.ID(),
		IsRadio:  false,
		Value: form.Value[[]form.Value[bool]]{
			Key: form.Key{
				Label: "Sources",
				Key:   "search-source-types",
			},
			Value: []form.Value[bool]{
				form.Value[bool]{
					Key: form.Key{
						Label: "Anime",
						Key:   "series_anime",
					},
					Value: true,
				},
				form.Value[bool]{
					Key: form.Key{
						Label: "Anime Movie",
						Key:   "movie_anime",
					},
					Value: true,
				},
				form.Value[bool]{
					Key: form.Key{
						Label: "Manga",
						Key:   "book_manga",
					},
					Value: true,
				},
				form.Value[bool]{
					Key: form.Key{
						Label: "Light Novel",
						Key:   "book_light_novel",
					},
					Value: true,
				},
			},
		},
	})
	n.options = checkbox_group.New(checkbox_group.O{
		Prefix: "search-options",
		ParentID: n.ID(),
		IsRadio: false,
		Value: form.Value[[]form.Value[bool]]{
			Key: form.Key{
				Label: "Options",
				Key: "search-options",
			},
			Value: []form.Value[bool]{
				form.Value[bool]{
					Key: form.Key{
						Label: "NSFW",
						Key: "options_nsfw",
					},
					Value: false,
				},
			},
		},
	})
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.input.Init(),
		n.apis.Init(),
		n.source_types.Init(),
		n.options.Init(),
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

	for _, m := range []tea.Model{
		n.input,
		n.apis,
		n.source_types,
		n.options,
	} {
		_, c = m.Update(msg)
		cmds = append(cmds, c)
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return n.column.RenderOrDie(
		lipgloss.JoinVertical(
			lipgloss.Left,
			n.column.Style().Render(n.input.View()),
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				n.apis.View(),
				lipgloss.NewStyle().MarginLeft(1).Render(
					n.source_types.View(),
				),
				lipgloss.NewStyle().MarginLeft(1).Render(
					n.options.View(),
				),
			),
		),
	)
}
