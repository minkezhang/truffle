package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/checkbox"
	"github.com/minkezhang/truffle/tui/component/checkbox_group"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/org"
	"github.com/minkezhang/truffle/tui/component/textarea"
	"github.com/minkezhang/truffle/tui/component/textinput"
	"github.com/minkezhang/truffle/tui/component/viewport"
)

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

const max_width = 100

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	c := column.New(max_width)

	pg := page{
		children: []tea.Model{
			org.New(c),
			textinput.New(textinput.O{
				Prefix:      "test textinput",
				ParentID:    "",
				Width:       max_width,
				Placeholder: "this is some text placeholder",
				Prompt:      "> ",
				Value:       "Frieren",
			}),
			textarea.New(textarea.O{
				Prefix:      "test textarea",
				ParentID:    "",
				Width:       max_width,
				Height:      10,
				Placeholder: "this is some textarea placeholder",
				Value:       "This is a synopsis",
				Label:       "Synopsis",
			}),
			checkbox.New(checkbox.O{
				Prefix:     "api mal",
				ParentID:   "",
				Label:      "MAL",
				Value:      "mal",
				IsSelected: false,
				IsRadio:    false,
			}),
			checkbox.New(checkbox.O{
				Prefix:     "type book",
				ParentID:   "",
				Label:      "book",
				Value:      "book",
				IsSelected: false,
				IsRadio:    true,
			}),
			checkbox_group.New(checkbox_group.O{
				Prefix:   "inputgroup-api",
				ParentID: "",
				IsRadio:  false,
				Values: []checkbox_group.Value{
					{
						Label:      "MAL",
						Value:      "mal",
						IsSelected: true,
					},
					{
						Label:      "Truffle",
						Value:      "truffle",
						IsSelected: true,
					},
					{
						Label:      "OMDB",
						Value:      "omdb",
						IsSelected: false,
					},
				},
			}),
			checkbox_group.New(checkbox_group.O{
				Prefix:   "inputgroup-type",
				ParentID: "",
				IsRadio:  true,
				Values: []checkbox_group.Value{
					{
						Label:      "Book",
						Value:      "book",
						IsSelected: true,
					},
					{
						Label:      "Anime",
						Value:      "anime",
						IsSelected: true,
					},
					{
						Label:      "Movie",
						Value:      "movie",
						IsSelected: false,
					},
				},
			}),
		},
	}

	rt := root{}
	rt.viewport = viewport.New(viewport.O{
		Prefix:   "viewport",
		ParentID: "",
		Column:   c,
		Node:     pg,
	})
	p := tea.NewProgram(rt, tea.WithAltScreen(), tea.WithMouseAllMotion())
	errors.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v\n", err)
		os.Exit(1)
	}
}
