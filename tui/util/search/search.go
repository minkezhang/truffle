package util_search

import (
	"context"
	"errors"

	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/db"
	"github.com/minkezhang/truffle/tui/util/node"
	"github.com/minkezhang/truffle/tui/util/node/virtual"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func Search(
	ctx context.Context,
	db *db.DB, query string,
	opts map[epb.SourceAPI][]option.O,
	types []epb.SourceType,
) ([]util_node.N, error) {
	var results []util_node.N
	var _errors []error
	node_ids := map[string]bool{}
	_types := map[epb.SourceType]bool{}

	for _, t := range types {
		_types[t] = true
	}

	sources, err := db.Search(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	for _, s := range sources {
		if !_types[s.Header().Type()] {
			continue
		}
		if s.NodeID() == "" {
			results = append(results, virtual.Make(s.PB()))
		} else if !node_ids[s.NodeID()] {
			node_ids[s.NodeID()] = true

			n, err := db.GetNode(ctx, node.Make(
				&dpb.Node{
					Header: &dpb.NodeHeader{
						Id:   s.NodeID(),
						Type: s.Header().Type(),
					},
				},
			).Header(), option.Remote(false))
			if err != nil {
				_errors = append(_errors, err)
			} else {
				results = append(results, n)
			}
		}
	}
	return results, errors.Join(_errors...)
}
