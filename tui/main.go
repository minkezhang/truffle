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
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/log"
	"github.com/minkezhang/truffle/tui/component/textarea"
	"github.com/minkezhang/truffle/tui/component/textinput"
)

type root struct {
	directory      tea.Model
	textinput      tea.Model
	errors         tea.Model
	log            tea.Model
	textarea       tea.Model
	checkbox       tea.Model
	radio          tea.Model
	checkbox_group tea.Model
	radio_group    tea.Model
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		tea.Sequence( // preserve tab order
			r.directory.Init(),
			r.textinput.Init(),
			r.errors.Init(),
			r.log.Init(),
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
		r.directory,
		r.textinput,
		r.errors,
		r.log,
		r.textarea,
		r.checkbox,
		r.radio,
		r.checkbox_group,
		r.radio_group,
	} {
		_, c = n.Update(msg)
		cmds = append(cmds, c)
	}
	return r, tea.Batch(cmds...)
}

func (r root) View() string {
	return zone.Scan(
		lipgloss.JoinVertical(
			lipgloss.Left,
			r.directory.View(),
			r.textinput.View(),
			r.log.View(),
			r.textarea.View(),
			r.checkbox.View(),
			r.radio.View(),
			r.checkbox_group.View(),
			r.radio_group.View(),
		),
	)
}

const max_width = 100

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	c := column.New(max_width)
	c.SetBorder(lipgloss.NormalBorder(), true, false, true, false)

	p := tea.NewProgram(
		root{
			directory: focusable.New(),
			textinput: textinput.New(textinput.O{
				Prefix:      "test textinput",
				ParentID:    "",
				Width:       max_width,
				Placeholder: "this is some text placeholder",
				Prompt:      "> ",
				Value:       "Frieren",
			}),
			errors: &errors.Node{},
			log:    log.New("", c),
			textarea: textarea.New(textarea.O{
				Prefix:      "test textarea",
				ParentID:    "",
				Width:       max_width,
				Height:      10,
				Placeholder: "this is some textarea placeholder",
				Value:       "This is a synopsis",
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
				Values: []checkbox_group.V{
					{
						L:          "MAL",
						V:          "mal",
						IsSelected: true,
					},
					{
						L:          "Truffle",
						V:          "truffle",
						IsSelected: true,
					},
					{
						L:          "OMDB",
						V:          "omdb",
						IsSelected: false,
					},
				},
			}),
			radio_group: checkbox_group.New(checkbox_group.O{
				Prefix:   "inputgroup-type",
				ParentID: "",
				IsRadio:  true,
				Values: []checkbox_group.V{
					{
						L:          "Book",
						V:          "book",
						IsSelected: true,
					},
					{
						L:          "Anime",
						V:          "anime",
						IsSelected: true,
					},
					{
						L:          "Movie",
						V:          "movie",
						IsSelected: false,
					},
				},
			}),
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
