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
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/org"
	"github.com/minkezhang/truffle/tui/component/search"
	"github.com/minkezhang/truffle/tui/component/table"
	"github.com/minkezhang/truffle/tui/component/viewport"
	"github.com/minkezhang/truffle/tui/util/form"
	"github.com/minkezhang/truffle/tui/util/node"
	"github.com/minkezhang/truffle/tui/util/node/virtual"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
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
		nil,
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
				Data: []util_node.N{
					node.Make(
						&dpb.Node{
							Header: &dpb.NodeHeader{
								Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
								Id:   "",
							},
						},
					).WithSources(
						[]source.S{
							source.Make(
								&dpb.Source{
									Header: &dpb.SourceHeader{
										Api:  epb.SourceAPI_SOURCE_API_MAL,
										Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
										Id:   "523",
									},
									Titles: []*dpb.Title{
										&dpb.Title{
											Title:        "Meitantei Conan",
											Localization: "",
										},
										&dpb.Title{
											Title:        "Case Closed",
											Localization: "en",
										},
									},
									Score:      81,
									PreviewUrl: "https://cdn.myanimelist.net/images/anime/7/75199l.jpg",
								},
							),
						},
					),
					node.Make(
						&dpb.Node{
							Header: &dpb.NodeHeader{
								Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
								Id:   "",
							},
						},
					).WithSources(
						[]source.S{
							source.Make(
								&dpb.Source{
									Header: &dpb.SourceHeader{
										Api:  epb.SourceAPI_SOURCE_API_MAL,
										Type: epb.SourceType_SOURCE_TYPE_SERIES_ANIME,
										Id:   "52991",
									},
									Titles: []*dpb.Title{
										&dpb.Title{
											Title:        "Sousou no Frieren",
											Localization: "ja",
										},
										&dpb.Title{
											Title:        "Frieren: Beyond Journey's End",
											Localization: "en",
										},
									},
									Score:      92,
									PreviewUrl: "https://cdn.myanimelist.net/images/anime/1015/138006l.jpg",
								},
							),
						},
					),
				},
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

func do_search(msg search.SubmitSearchMessage) tea.Cmd {
	return func() tea.Msg {
		types := map[epb.SourceType]bool{}
		for _, t := range msg.SourceTypes.Value {
			if t.Value {
				if v, ok := epb.SourceType_value[t.Key.Key]; ok {
					types[epb.SourceType(v)] = true
				}
			}
		}

		_opts := []option.O{option.Remote(true)}
		for _, o := range msg.Options.Value {
			if o.Value {
				_opts = append(_opts, map[string]option.O{
					"options-nsfw": option.NSFW(true),
				}[o.Key.Key])
			}
		}

		opts := map[epb.SourceAPI][]option.O{}
		for _, api := range msg.APIs.Value {
			if api.Value {
				if v, ok := epb.SourceAPI_value[api.Key.Key]; ok {
					opts[epb.SourceAPI(v)] = append([]option.O{}, _opts...)
				}
			}
		}
		results, err := _db.Search(
			context.Background(),
			msg.Query.Value,
			opts,
		)
		if err != nil {
			return errors.ToLogMessage(
				errors.LevelWarn,
				err.Error(),
			)
		}

		var nodes []util_node.N
		for _, r := range results {
			if types[r.Header().Type()] {
				nodes = append(
					nodes,
					virtual.Make(r.PB()),
				)
			}
		}

		return search_result_message{
			results: nodes,
		}
	}
}

type search_result_message struct {
	results []util_node.N
}

func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	switch msg := msg.(type) {
	case search.SubmitSearchMessage:
		cmds = append(cmds, do_search(msg))
	case search_result_message:
		cmds = append(cmds, p.children[2].(*table.Node).SetValues(msg.results))
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
}

func (r root) Init() tea.Cmd {
	return tea.Sequence(
		r.errors.Init(),
		r.viewport.Init(),
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
	}
	p := tea.NewProgram(rt, tea.WithAltScreen(), tea.WithMouseCellMotion())
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
