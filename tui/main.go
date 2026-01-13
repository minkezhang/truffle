package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/checkbox"
	"github.com/minkezhang/truffle/tui/component/checkbox_group"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/org"
	"github.com/minkezhang/truffle/tui/component/textarea"
	"github.com/minkezhang/truffle/tui/component/textinput"
)

type root struct {
	org            tea.Model
	textinput      tea.Model
	textarea       tea.Model
	checkbox       tea.Model
	radio          tea.Model
	checkbox_group tea.Model
	radio_group    tea.Model

	viewport viewport.Model
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		tea.Sequence( // preserve tab order
			r.org.Init(),
			r.textinput.Init(),
			r.textarea.Init(),
			r.checkbox.Init(),
			r.radio.Init(),
			r.checkbox_group.Init(),
			r.radio_group.Init(),
		),
		func() tea.Msg {
			return types.FocusMessage{
				ID: r.textinput.(base.Identifiable).ID(),
			}
		},
	)
}

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

	for _, n := range []tea.Model{
		r.org,
		r.textinput,
		r.textarea,
		r.checkbox,
		r.radio,
		r.checkbox_group,
		r.radio_group,
	} {
		_, c = n.Update(msg)
		cmds = append(cmds, c)
	}
	r.viewport, c = r.viewport.Update(msg)
	cmds = append(cmds, c)
	r.viewport.SetContent(r.view())
	return r, tea.Batch(cmds...)
}

func (r root) View() string { return zone.Scan(r.viewport.View()) }
func (r root) view() string {
	// return zone.Scan(
	return lipgloss.JoinVertical(
		lipgloss.Left,
		r.textinput.View(),
		r.textarea.View(),
		r.checkbox.View(),
		r.radio.View(),
		r.checkbox_group.View(),
		r.radio_group.View(),
		r.org.View(),
	//	),
	)
}

const max_width = 100

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	c := column.New(max_width)

	p := tea.NewProgram(
		root{
			org: org.New(c),
			textinput: textinput.New(textinput.O{
				Prefix:      "test textinput",
				ParentID:    "",
				Width:       max_width,
				Placeholder: "this is some text placeholder",
				Prompt:      "> ",
				Value:       "Frieren",
			}),
			textarea: textarea.New(textarea.O{
				Prefix:      "test textarea",
				ParentID:    "",
				Width:       max_width,
				Height:      10,
				Placeholder: "this is some textarea placeholder",
				Value:       "This is a synopsis",
				Label:       "Synopsis",
			}),
			checkbox: checkbox.New(checkbox.O{
				Prefix:     "api mal",
				ParentID:   "",
				Label:      "MAL",
				Value:      "mal",
				IsSelected: false,
				IsRadio:    false,
			}),
			radio: checkbox.New(checkbox.O{
				Prefix:     "type book",
				ParentID:   "",
				Label:      "book",
				Value:      "book",
				IsSelected: false,
				IsRadio:    true,
			}),
			checkbox_group: checkbox_group.New(checkbox_group.O{
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
			radio_group: checkbox_group.New(checkbox_group.O{
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
			viewport: viewport.New(c.Content(), 40),
		},
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)
	errors.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v\n", err)
		os.Exit(1)
	}
}
