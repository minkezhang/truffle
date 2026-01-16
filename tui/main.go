package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/button"
	"github.com/minkezhang/truffle/tui/component/checkbox"
	"github.com/minkezhang/truffle/tui/component/checkbox_group"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/org"
	"github.com/minkezhang/truffle/tui/component/search"
	"github.com/minkezhang/truffle/tui/component/table"
	"github.com/minkezhang/truffle/tui/component/textarea"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/component/viewport"
	"github.com/minkezhang/truffle/tui/util/form"
)

func make_page(c *column.C) page {
	return page{
		children: []tea.Model{
			org.New(c),
			textinput.New(textinput.O{
				Prefix:      "test textinput",
				ParentID:    "",
				Width:       c.Content(),
				Placeholder: "this is some text placeholder",
				Prompt:      "> ",
				Value: form.Value[string]{
					Key: form.Key{
						Label: "Search",
						Key:   "form-key",
					},
					Value: "Frieren",
				},
			}),
			textarea.New(textarea.O{
				Prefix:      "test textarea",
				ParentID:    "",
				Width:       c.Content(),
				Height:      10,
				Placeholder: "this is some textarea placeholder",
				Value: form.Value[string]{
					Key: form.Key{
						Label: "Synopsis",
						Key:   "source-synopsis",
					},
					Value: "This is a synopsis",
				},
			}),
			checkbox.New(checkbox.O{
				Prefix:   "api mal",
				ParentID: "",
				Value: form.Value[bool]{
					Key: form.Key{
						Label: "MAL",
						Key:   "mal",
					},
					Value: false,
				},
				IsRadio: false,
			}),
			checkbox.New(checkbox.O{
				Prefix:   "type book",
				ParentID: "",
				Value: form.Value[bool]{
					Key: form.Key{
						Label: "Book",
						Key:   "book",
					},
					Value: false,
				},
				IsRadio: true,
			}),
			checkbox_group.New(checkbox_group.O{
				Prefix:   "inputgroup-api",
				ParentID: "",
				IsRadio:  false,
				Value: form.Value[[]form.Value[bool]]{
					Key: form.Key{
						Label: "APIs",
						Key:   "apis",
					},
					Value: []form.Value[bool]{
						form.Value[bool]{
							Key: form.Key{
								Label: "MAL",
								Key:   "mal",
							},
							Value: true,
						},
						form.Value[bool]{
							Key: form.Key{
								Label: "Truffle",
								Key:   "truffle",
							},
							Value: true,
						},
						form.Value[bool]{
							Key: form.Key{
								Label: "OMDB",
								Key:   "omdb",
							},
							Value: false,
						},
					},
				},
			}),
			checkbox_group.New(checkbox_group.O{
				Prefix:   "inputgroup-source-type",
				ParentID: "",
				IsRadio:  true,
				Value: form.Value[[]form.Value[bool]]{
					Key: form.Key{
						Label: "Source Type",
						Key:   "source-type",
					},
					Value: []form.Value[bool]{
						form.Value[bool]{
							Key: form.Key{
								Label: "BOOK",
								Key:   "book",
							},
							Value: true,
						},
						form.Value[bool]{
							Key: form.Key{
								Label: "ANIME",
								Key:   "anime",
							},
							Value: true,
						},
						form.Value[bool]{
							Key: form.Key{
								Label: "MOVIE",
								Key:   "movie",
							},
							Value: false,
						},
					},
				},
			}),
			button.New(button.O{
				ParentID: "",
				Key: form.Key{
					Label: "Test",
					Key:   "test",
				},
			}),
			search.New(search.O{
				Column: c,
			}),
			table.New(table.O{
				Prefix: "table-view",
			}),
		},
	}
}

type page struct {
	children []tea.Model
}

func (p page) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range p.children {
		cmds = append(cmds, c.Init())
	}
	return tea.Sequence(cmds...)
}

func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	for i := range p.children {
		p.children[i], c = p.children[i].Update(msg)
		cmds = append(cmds, c)
	}
	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	var parts []string
	for _, c := range p.children {
		parts = append(parts, c.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

type root struct {
	viewport tea.Model
}

func (r root) Init() tea.Cmd { return r.viewport.Init() }

func (r root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyCtrlD:
			cmds = append(cmds, tea.Quit)
		case tea.KeyCtrlZ:
			cmds = append(cmds, tea.Suspend)
		}
	case tea.WindowSizeMsg:
		cmds = append(cmds, tea.ClearScreen)
	}

	r.viewport, c = r.viewport.Update(msg)
	cmds = append(cmds, c)

	return r, tea.Batch(cmds...)
}

func (r root) View() string { return zone.Scan(r.viewport.View()) }

const max_width = 150

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	c := column.New(max_width)
	rt := root{
		viewport: viewport.New(viewport.O{
			Prefix:   "viewport",
			ParentID: "",
			Column:   c,
			Node:     make_page(c.WithWidth(c.Width() - 2)),
		}),
	}
	p := tea.NewProgram(rt, tea.WithAltScreen(), tea.WithMouseCellMotion())
	errors.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v\n", err)
		os.Exit(1)
	}
}
