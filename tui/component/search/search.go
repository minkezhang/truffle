package search

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/button"
	"github.com/minkezhang/truffle/tui/component/checkbox_group"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"

	directory "github.com/minkezhang/truffle/tui/component/directory/focusable"
)

type Node struct {
	*focusable.Node

	column *column.C
	input  *textinput.Node

	apis          *checkbox_group.Node
	source_types  *checkbox_group.Node
	options       *checkbox_group.Node
	submit_button *button.Node
	is_expanded   bool
	clickable     *clickable.Node
}

type O struct {
	ParentID string
	Column   *column.C
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New("search", o.ParentID, 1),
		column: o.Column,
	}
	n.input = textinput.New(textinput.O{
		Prefix:      "search-textinput",
		ParentID:    n.ID(),
		Width:       o.Column.Content() - 4,
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
				Label: "Media",
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
		Prefix:   "search-options",
		ParentID: n.ID(),
		IsRadio:  false,
		Value: form.Value[[]form.Value[bool]]{
			Key: form.Key{
				Label: "Options",
				Key:   "search-options",
			},
			Value: []form.Value[bool]{
				form.Value[bool]{
					Key: form.Key{
						Label: "NSFW",
						Key:   "options_nsfw",
					},
					Value: false,
				},
			},
		},
	})
	n.submit_button = button.New(button.O{
		ParentID: n.ID(),
		Key: form.Key{
			Label: "Search",
			Key:   "search-submit",
		},
	})
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	cmds := []tea.Cmd{
		n.input.Init(),
		n.clickable.Init(),
		n.apis.Init(),
		n.source_types.Init(),
		n.options.Init(),
		n.submit_button.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	}
	for _, c := range []directory.Node{
		n.apis,
		n.source_types,
		n.options,
		n.submit_button,
	} {
		cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
	}
	return tea.Sequence(cmds...)
}

type SubmitSearchMessage struct {
	ID          string
	Query       form.Value[string]
	APIs        form.Value[[]form.Value[bool]]
	SourceTypes form.Value[[]form.Value[bool]]
	Options     form.Value[[]form.Value[bool]]
}

func (n *Node) submit() tea.Cmd {
	m := SubmitSearchMessage{
		ID:          n.ID(),
		Query:       n.input.Value(),
		APIs:        n.apis.Value(),
		SourceTypes: n.source_types.Value(),
		Options:     n.options.Value(),
	}

	// Ignore blank queries.
	if m.Query.Value == "" {
		return nil
	}

	return tea.Sequence(
		n.input.SetValue(""),
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: submitting search query %v", n.ID(), m),
			)
		},
		func() tea.Msg {
			return m
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
		n.submit_button,
		n.clickable,
	} {
		_, c = m.Update(msg)
		cmds = append(cmds, c)
	}

	switch msg := msg.(type) {
	case button.SubmitMessage:
		if msg.ID == n.submit_button.ID() {
			cmds = append(cmds, n.submit())
		}
	case tea.KeyMsg:
		if n.input.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter {
				cmds = append(cmds, n.submit())
			}
		}
		if n.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				n.is_expanded = !n.is_expanded
				for _, c := range []directory.Node{
					n.apis,
					n.source_types,
					n.options,
					n.submit_button,
				} {
					cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
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
			for _, c := range []directory.Node{
				n.apis,
				n.source_types,
				n.options,
				n.submit_button,
			} {
				cmds = append(cmds, c.SetIsInvisible(!n.is_expanded))
			}
		}
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	parts := []string{
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			n.input.View(),
			lipgloss.NewStyle().PaddingLeft(1).Foreground(
				color_profile.UIForeground[n.FocusState()],
			).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
				color_profile.UIForeground[n.input.FocusState()],
			).Render(
				zone.Mark(
					n.clickable.ID(),
					map[bool]string{
						false: "(+)",
						true:  "(-)",
					}[n.is_expanded],
				),
			),
		),
	}
	if n.is_expanded {
		parts = append(parts,
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
			n.submit_button.View(),
		)
	}
	return n.column.RenderOrDie(lipgloss.JoinVertical(lipgloss.Left, parts...))
}
