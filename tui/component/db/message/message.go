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

type SearchRequestMessage struct {
	ID          string
	Query       form.Value[string]
	APIs        form.Value[map[epb.SourceAPI]bool]
	SourceTypes form.Value[map[epb.SourceType]bool]
	Options     form.Value[[]option.O]
}

type SearchResponseMessage struct {
	ID      string
	Results []util_node.N
}

type GetNodeRequestMessage struct {
	ID    string
	Value form.Value[util_node.N]
}

type PutRequestMessage struct {
	ID    string
	Value form.Value[source.S]
}

type PutResponseMessage struct {
	ID          string
	Node        node.N
	SourceIndex int
}
