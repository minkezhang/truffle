package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
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
	directory tea.Model
	textinput tea.Model
	errors    tea.Model
	log       tea.Model
	textarea  tea.Model
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		r.directory.Init(),
		r.textinput.Init(),
		r.errors.Init(),
		r.log.Init(),
		r.textarea.Init(),
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
			r.log.View(),
			r.textinput.View(),
			r.textarea.View(),
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
