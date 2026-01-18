package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/org"
	"github.com/minkezhang/truffle/tui/component/search"
	"github.com/minkezhang/truffle/tui/component/table"
	"github.com/minkezhang/truffle/tui/component/viewport"
	"github.com/minkezhang/truffle/tui/util/form"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	component_db "github.com/minkezhang/truffle/tui/component/db"
)

const (
	cache       = "./.build/"
	MALClientID = "6114d00ca681b7701d1e15fe11a4987e"
)

var (
	_db = db.New(
		context.Background(),
		&cpb.Config{
			Mal: &cpb.MAL{
				ClientId:   MALClientID,
				MaxResults: 10,
			},
			Truffle: &cpb.Truffle{},
		},
		&dpb.Database{
			Nodes: []*dpb.Node{
				&dpb.Node{
					Header: &dpb.NodeHeader{
						Id:   "frieren-anime",
						Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
					},
				},
			},
			Sources: []*dpb.Source{
				&dpb.Source{
					NodeId: "frieren-anime",
					Header: &dpb.SourceHeader{
						Api:  epb.SourceAPI_SOURCE_API_MAL,
						Id:   "52991",
						Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
					},
				},
				&dpb.Source{
					NodeId: "frieren-anime",
					Header: &dpb.SourceHeader{
						Api:  epb.SourceAPI_SOURCE_API_TRUFFLE,
						Id:   "frieren-anime-truffle",
						Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
					},
					Titles: []*dpb.Title{
						&dpb.Title{Title: "Frieren", Localization: "en"},
					},
					Synopsis: "An anime series",
				},
			},
		},
	)
)

func make_page(c *column.C) page {
	return page{
		children: []tea.Model{
			org.New(c),
			search.New(search.O{
				Column: c,
			}),
			table.New(table.O{
				Prefix: "table-view",
				Column: c,
				Key: form.Key{
					Label: "Results",
					Key:   "search-results",
				},
				CacheDirectory: cache,
			}),
		},
	}
}

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

	switch msg := msg.(type) {
	case component_db.SearchResultMessage:
		cmds = append(cmds, p.children[2].(*table.Node).SetValues(msg.Results))
	}

	for i := range p.children {
		p.children[i], c = p.children[i].Update(msg)
		cmds = append(cmds, c)
	}
	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	var parts []string
	for _, c := range p.children {
		if t, ok := c.(*table.Node); ok { // Only display table if there are results
			if !t.IsInvisible() {
				parts = append(parts, c.View())
			}
		} else {
			parts = append(parts, c.View())
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

type root struct {
	errors   tea.Model
	viewport tea.Model
	db       tea.Model
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		r.errors.Init(),
		r.viewport.Init(),
		r.db.Init(),
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
		r.errors,
		r.viewport,
		r.db,
	} {
		_, c = n.Update(msg)
		cmds = append(cmds, c)
	}

	return r, tea.Batch(cmds...)
}

func (r root) View() string { return zone.Scan(r.viewport.View()) }

const max_width = 175

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if _, err := tea.LogToFile(filepath.Join(cache, "debug.log"), ""); err != nil {
		fmt.Printf("cannot open error log")
		os.Exit(1)
	}

	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	c := column.New(max_width)
	rt := root{
		errors: &errors.Node{},
		viewport: viewport.New(viewport.O{
			Prefix:   "viewport",
			ParentID: "",
			Column:   c,
			Node:     make_page(c.WithWidth(c.Width() - 2)),
		}),
		db: component_db.New(context.Background(), component_db.O{
			DB: _db,
		}),
	}
	p := tea.NewProgram(
		rt,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithoutCatchPanics(),
	)
	errors.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v\n", err)
		os.Exit(1)
	}

	if err := rt.errors.(*errors.Node).Error(); err != nil {
		fmt.Printf("Error() returned unexpected error: %v\n", err)
		os.Exit(1)
	}
}
