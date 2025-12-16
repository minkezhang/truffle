package root

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/client"
	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle-api/db/query"
	"github.com/minkezhang/truffle/tui/util/grid"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	node_ui "github.com/minkezhang/truffle/tui/component/node"
	source_list_ui "github.com/minkezhang/truffle/tui/component/node/source_list"
)

type O struct {
	CacheDirectory string
}

type M struct {
	directory string
	node      tea.Model
	debug     *source_list_ui.M
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
		directory: o.CacheDirectory,
		node: node_ui.New(node_ui.O{
			Layout:         grid.L{Content: 120, Width: 120},
			CacheDirectory: o.CacheDirectory,
			Node:           ns[0],
		}),
		debug: source_list_ui.New(source_list_ui.O{
			Layout: grid.L{Content: 50, Width: 50},
			Values: []source_list_ui.V{
				source_list_ui.V{
					API: epb.API_API_VIRTUAL,
					ID:  "",
				},
				source_list_ui.V{
					API: epb.API_API_TRUFFLE,
					ID:  "abc",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "123",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "4",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "5",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "6",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "7",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "8",
				},
				source_list_ui.V{
					API: epb.API_API_MAL,
					ID:  "9",
				},
			},
		}),
	}
}

func (m *M) Init() tea.Cmd { return m.node.Init() }

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
	}

	_, c := m.debug.Update(msg)
	cmds = append(cmds, c)

	_, c = m.node.Update(msg)
	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)
}

func (m *M) View() string {
	return zone.Scan(m.debug.View())
	// return zone.Scan(m.node.View())
}
