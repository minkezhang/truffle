package root

import (
	"context"

	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/client/query"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"

	tuibook "github.com/minkezhang/truffle/tui/component/metadata/book"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type M struct {
	metadata *tuibook.M
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
		metadata: tuibook.New(tuibook.O{
			Book: a.Metadata().(*book.M),
		}),
	}
}
