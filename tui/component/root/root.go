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
	"github.com/minkezhang/truffle/tui/component/util/overlay"
	"github.com/minkezhang/truffle/tui/util/grid"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	node_card_ui "github.com/minkezhang/truffle/tui/component/node/card"
	node_full_ui "github.com/minkezhang/truffle/tui/component/node/full"
	search_ui "github.com/minkezhang/truffle/tui/component/search"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type ViewMode int

const (
	ViewModeNone ViewMode = iota
	ViewModeCard
	ViewModeFull
)

var (
	column = grid.C{Content: 120}
)

type O struct {
	CacheDirectory string
}

type M struct {
	truffle *db.DB

	directory string
	mode      ViewMode
	full      tea.Model
	search    tea.Model
	card      tea.Model

	// e is the global error handler
	e tea.Model
}

func New(o O) *M {
	truffle, err := db.New(context.Background(), db.O{
		Clients: []client.C{
			mal.New(mal.O{
				ClientID:         "6114d00ca681b7701d1e15fe11a4987e",
				PopularityCutoff: 10000,
				MaxResults:       20,
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
		card: node_card_ui.Make(node_card_ui.O{
			O: model_ui.O{
				Column: column.Grid(4).Column(2).WithPadding(1),
			},
			Nodes: nil,
		}),
		full: node_full_ui.Make(node_full_ui.O{
			O: model_ui.O{
				Column: column,
			},
			CacheDirectory: o.CacheDirectory,
			Node:           ns[0],
		}),

		search: search_ui.New(search_ui.O{
			O: model_ui.O{
				Column: column,
			},
		}),
		e:    model_ui.E,
		mode: ViewModeFull,
	}
}

func (m *M) Init() tea.Cmd {
	return tea.Batch(
		tea.Sequence(
			m.search.Init(),
			model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{ID: m.search.(model_ui.I).ID()},
			}),
		),
		m.full.Init(),
		m.card.Init(),
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
		case tea.KeyCtrlL:
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{ID: m.search.(model_ui.I).ID()},
			}))
		}
	case tea.WindowSizeMsg: // Clear buffer
		return m, tea.ClearScreen
	case search_ui.QueryMsg:
		cmds = append(cmds, model_ui.ToCommand(m.query(msg)))
	case model_ui.BlurMsg:
		switch id := msg.ID; id {
		case m.full.(model_ui.I).ID():
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{ID: m.search.(model_ui.I).ID()},
			}))
		case m.search.(model_ui.I).ID():
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{ID: m.full.(model_ui.I).ID()},
			}))
		case m.card.(model_ui.I).ID():
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{ID: m.search.(model_ui.I).ID()},
			}))
		}
	case node_card_ui.SelectMsg:
		m.mode = ViewModeFull
		m.full = node_full_ui.Make(node_full_ui.O{
			O: model_ui.O{
				Column: column,
			},
			CacheDirectory: m.directory,
			Node:           (*node.N)(msg),
		})
		cmds = append(
			cmds,
			model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{ID: m.full.(model_ui.I).ID()},
			}),
			m.full.Init(),
		)
	case ResponseMsg:
		m.mode = ViewModeCard
		m.card = node_card_ui.Make(node_card_ui.O{
			O: model_ui.O{
				Column: column.Grid(4).Column(2).WithPadding(1),
			},
			Nodes: []*node.N(msg),
		})
		cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
			BaseMsg: model_ui.BaseMsg{ID: m.card.(model_ui.I).ID()},
		}))
	}

	var c tea.Cmd

	m.full, c = m.full.Update(msg)

	cmds = append(cmds, c)

	m.e, c = m.e.Update(msg)

	cmds = append(cmds, c)

	m.search, c = m.search.Update(msg)

	cmds = append(cmds, c)

	m.card, c = m.card.Update(msg)

	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)
}

type ResponseMsg []*node.N

func (m *M) query(msg search_ui.QueryMsg) tea.Msg {
	ns, err := m.truffle.Query(context.Background(), query.New(query.O{
		APIs:      []epb.API{epb.API_API_MAL},
		AtomTypes: []epb.Type{epb.Type_TYPE_BOOK},
		Title:     string(msg),
	}))
	if err != nil {
		return model_ui.ErrorMsg(err)
	}
	return ResponseMsg(ns)
}

func (m *M) body() string {
	switch v := m.mode; v {
	case ViewModeFull:
		return overlay.M{FG: m.full, BG: m.card}.View()
	case ViewModeCard:
		return overlay.M{BG: m.full, FG: m.card}.View()
	}
	return ""
}

func (m *M) View() string {
	return zone.Scan(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.search.View(),
			m.body(),
		),
	)
}
