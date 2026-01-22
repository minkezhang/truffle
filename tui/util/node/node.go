package util_node

import (
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
)

type N interface {
	Header() node.H
	Sources() []source.S
	PB() *dpb.Node
	Virtual() (source.S, error)
}
