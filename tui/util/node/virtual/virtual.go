package virtual

import (
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
)

type N struct {
	pb *dpb.Source
}

func Make(pb *dpb.Source) N {
	return N{pb: pb}
}

func (n N) Header() node.H {
	return node.Make(&dpb.Node{
		Header: &dpb.NodeHeader{
			Id:   "",
			Type: n.pb.GetHeader().GetType(),
		},
	}).Header()
}

func (n N) Sources() []source.S { return []source.S{source.Make(n.pb)} }

func (n N) PB() *dpb.Node {
	return node.Make(&dpb.Node{
		Header: n.Header().PB(),
	}).PB()
}

func (n N) Virtual() (source.S, error) { return source.Make(n.pb), nil }
