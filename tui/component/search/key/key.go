package key

import (
	"github.com/minkezhang/truffle/tui/util/form"
)

var (
	// Search is the exported search bar key sent to the DB as part of the
	// request. This key is returned in the response, which may be consumed
	// by other UI nodes. Exporting this key allows other nodes to know the
	// originating UI node.
	Search = form.Key{Key: "search-query"}
)
