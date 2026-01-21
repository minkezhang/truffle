// Package message contains request / response messages for interacting with the
// Truffle DB. Requests may be issued by any UI node, but reponses are issued by
// the dedicated DB shim.
package message

import (
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/util/form"
	"github.com/minkezhang/truffle/tui/util/node"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type SearchRequestBody struct {
	Query       string
	APIs        map[epb.SourceAPI]bool
	SourceTypes map[epb.SourceType]bool
	Options     []option.O
}

type SearchRequestMessage struct {
	ID   string
	Body form.Value[SearchRequestBody]
}

type SearchResponseMessage struct {
	ID   string
	Body form.Value[[]util_node.N]
}

type GetNodeRequestMessage struct {
	ID   string
	Body form.Value[util_node.N]
}

type GetNodeResponseMessage struct {
	ID   string
	Body form.Value[util_node.N]
}

type PutRequestMessage struct {
	ID   string
	Body form.Value[source.S]
}

type PutResponseBody struct {
	Node        node.N
	SourceIndex int
}

type PutResponseMessage struct {
	ID   string
	Body form.Value[PutResponseBody]
}
