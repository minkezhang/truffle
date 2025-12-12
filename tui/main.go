package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/client/query"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"
	"github.com/minkezhang/truffle/tui/component/util/debug"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	tuibook "github.com/minkezhang/truffle/tui/component/metadata/book"
)

type M struct {
	atom  *atom.A
	debug tea.Model
}

func New() *M {
	c := mal.New(mal.O{
		ClientID:         "6114d00ca681b7701d1e15fe11a4987e",
		PopularityCutoff: 10000,
		MaxResults:       2,
		NSFW:             true,
	})
	a, _ := c.Get(context.Background(), query.G{
		AtomType: epb.Type_TYPE_BOOK,
		ID:       "107562",
	})

	return &M{
		atom:  a,
		debug: debug.Init(),
	}
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.debug.Update(msg)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyCtrlZ:
			return m, tea.Suspend
		}
	}
	return m, nil
}

func (m *M) View() string {
	var s strings.Builder

	s.WriteString(m.debug.View())
	s.WriteString(
		tuibook.Init(tuibook.O{
			Book: m.atom.Metadata().(*book.M),
		}).View(),
	)
	return zone.Scan(s.String())
}

func main() {
	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v", err)
		os.Exit(1)
	}
}
