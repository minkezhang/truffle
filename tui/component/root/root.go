package root

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/client"
	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle-api/db/node"
	"github.com/minkezhang/truffle-api/db/query"
	"github.com/minkezhang/truffle/tui/util/grid"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	node_ui "github.com/minkezhang/truffle/tui/component/node"
	search_ui "github.com/minkezhang/truffle/tui/component/search"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	CacheDirectory string
}

type M struct {
	truffle *db.DB

	directory string
	node      tea.Model
	search    tea.Model

	// e is the global error handler
	e tea.Model
}

func New(o O) *M {
	truffle, err := db.New(context.Background(), db.O{
		Clients: []client.C{
			mal.New(mal.O{
				ClientID:         "6114d00ca681b7701d1e15fe11a4987e",
				PopularityCutoff: 10000,
				MaxResults:       2,
				NSFW:             true,
			}),
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("new error: %v", err))
		return nil
	}

	ns, err := truffle.Query(context.Background(), query.New(query.O{
		APIs:      []epb.API{epb.API_API_MAL},
		AtomTypes: []epb.Type{epb.Type_TYPE_BOOK},
		Title:     "The Apothecary Diaries",
	}))
	if err != nil {
		slog.Error(fmt.Sprintf("query error: %v", err))
		return nil
	}

	return &M{
		truffle:   truffle,
		directory: o.CacheDirectory,
		node: node_ui.Make(node_ui.O{
			O: model_ui.O{
				Column: grid.C{Content: 120},
			},
			CacheDirectory: o.CacheDirectory,
			Node:           ns[0],
		}),

		search: search_ui.New(search_ui.O{}),
		e:      model_ui.E,
	}
}

func (m *M) Init() tea.Cmd {
	return tea.Batch(
		m.node.Init(),
		m.search.Init(),
	)
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyCtrlZ:
			return m, tea.Suspend
		}
	case tea.WindowSizeMsg: // Clear buffer
		return m, tea.ClearScreen
	case search_ui.QueryMsg:
		cmds = append(cmds, func() tea.Msg { return m.query(msg) })
	case QueryResponseMsg:
		if len([]*node.N(msg)) > 0 {
			m.node = node_ui.Make(node_ui.O{
				O: model_ui.O{
					Column: grid.C{Content: 120},
				},
				CacheDirectory: m.directory,
				Node:           []*node.N(msg)[0],
			})
			cmds = append(cmds, m.node.Init())
		}
	}

	var c tea.Cmd

	m.node, c = m.node.Update(msg)

	cmds = append(cmds, c)

	m.e, c = m.e.Update(msg)

	cmds = append(cmds, c)

	m.search, c = m.search.Update(msg)

	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)
}

type QueryResponseMsg []*node.N

func (m *M) query(msg search_ui.QueryMsg) tea.Msg {
	ns, err := m.truffle.Query(context.Background(), query.New(query.O{
		APIs:      []epb.API{epb.API_API_MAL},
		AtomTypes: []epb.Type{epb.Type_TYPE_BOOK},
		Title:     string(msg),
	}))
	if err != nil {
		return model_ui.ErrorMsg(err)
	}
	return QueryResponseMsg(ns)
}

func (m *M) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.search.View(),
		zone.Scan(
			lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(
				m.node.View(),
			),
		),
	)
}
