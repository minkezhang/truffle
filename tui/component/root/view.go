package root

import (
	"strings"

	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"

	tuibook "github.com/minkezhang/truffle/tui/component/metadata/book"
)

func (m *M) View() string {
	var s strings.Builder

	if m.overlay {
		s.WriteString(m.debug.View())
	}
	s.WriteString(
		tuibook.Init(tuibook.O{
			Book: m.atom.Metadata().(*book.M),
		}).View(),
	)
	return zone.Scan(s.String())
}
