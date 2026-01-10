package renderable

import (
	"github.com/charmbracelet/bubbletea"

	"github.com/minkezhang/truffle/tui/component/directory"
)

type Node interface {
	directory.Identifiable

	// TODO(minkezhang)
}
