package root

import (
	"context"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/client/query"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle/tui/component/util/debug"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type M struct {
	atom    *atom.A
	debug   tea.Model
	overlay bool
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
