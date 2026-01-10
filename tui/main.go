package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/directory"
	"github.com/minkezhang/truffle/tui/component/directory/types"
	"github.com/minkezhang/truffle/tui/component/textinput"
)

type root struct {
	directory tea.Model
	input     tea.Model
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		r.directory.Init(),
		r.input.Init(),
		func() tea.Msg {
			return types.FocusMessage{
				ID: r.input.(types.Node).ID(),
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
	r.directory, c = r.directory.Update(msg)
	cmds = append(cmds, c)
	r.input, c = r.input.Update(msg)
	cmds = append(cmds, c)
	return r, tea.Batch(cmds...)
}

func (r root) View() string { return r.directory.View() + "\n" + zone.Scan(r.input.View()) }

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	p := tea.NewProgram(
		root{
			directory: directory.New(),
			input: textinput.New(textinput.O{
				Prefix:      "test input",
				ParentID:    "",
				Width:       50,
				Placeholder: "some text dim",
				Prompt:      "> ",
			}),
		},
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v\n", err)
		os.Exit(1)
	}
}
