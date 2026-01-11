package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/textinput"
)

type root struct {
	directory tea.Model
	input     tea.Model
	errors    *errors.Node
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		r.directory.Init(),
		r.input.Init(),
		r.errors.Init(),
		func() tea.Msg {
			return types.FocusMessage{
				ID: r.input.(base.Identifiable).ID(),
			}
		},
	)
}

func (r root) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyCtrlD:
			return r, tea.Quit
		case tea.KeyCtrlZ:
			return r, tea.Suspend
		}
	}
	var cmds []tea.Cmd
	var c tea.Cmd
	_, c = r.directory.Update(msg)
	cmds = append(cmds, c)
	_, c = r.input.Update(msg)
	cmds = append(cmds, c)
	_, c = r.errors.Update(msg)
	cmds = append(cmds, c)
	return r, tea.Batch(cmds...)
}

func (r root) View() string {
	return zone.Scan(
		r.directory.View() + "\n" + r.input.View(),
	)
}

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	p := tea.NewProgram(
		root{
			directory: focusable.New(),
			input: textinput.New(textinput.O{
				Prefix:      "test input",
				ParentID:    "",
				Width:       100,
				Placeholder: "this is some text placeholder",
				Prompt:      "> ",
				Value:       "Frieren",
			}),
			errors: &errors.Node{},
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
