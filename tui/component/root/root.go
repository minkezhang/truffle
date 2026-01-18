// Package root encapsulates all display-related logic for Truffle.
//
// root > viewport > display
package root

import (
	"context"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/display"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/viewport"

	component_db "github.com/minkezhang/truffle/tui/component/db"
)

const (
	max_width = 175
)

type O struct {
	DB             *db.DB
	CacheDirectory string
}

func Make(o O) Node {
	return Node{
		errors: &errors.Node{},
		viewport: viewport.New(viewport.O{
			Column: column.New(max_width),
			Node: display.Make(display.O{
				Column:         column.New(max_width - 2),
				CacheDirectory: o.CacheDirectory,
			}),
		}),
		db: component_db.New(
			context.Background(),
			component_db.O{
				DB: o.DB,
			},
		),
	}
}

type Node struct {
	errors   *errors.Node
	viewport *viewport.Node
	db       *component_db.Node
}

func (n Node) Init() tea.Cmd {
	return tea.Sequence(
		n.errors.Init(),
		n.viewport.Init(),
		n.db.Init(),
	)
}

func (n Node) Error() error { return n.errors.Error() }

func (n Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		n.errors,
		n.viewport,
		n.db,
	} {
		_, c = n.Update(msg)
		cmds = append(cmds, c)
	}

	return n, tea.Batch(cmds...)
}

func (n Node) View() string { return zone.Scan(n.viewport.View()) }
