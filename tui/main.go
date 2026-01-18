package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/root"

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

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	if _, err := tea.LogToFile(filepath.Join(cache, "debug.log"), ""); err != nil {
		fmt.Printf("cannot open error log")
		os.Exit(1)
	}

	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	n := root.Make(root.O{
		DB:             _db,
		CacheDirectory: cache,
	})
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
	p := tea.NewProgram(n, opts...)
	errors.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v\n", err)
		os.Exit(1)
	}

	if err := n.Error(); err != nil {
		fmt.Printf("Error() returned unexpected error: %v\n", err)
		os.Exit(1)
	}
}
