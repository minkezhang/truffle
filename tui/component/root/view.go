package root

import (
	"github.com/lrstanley/bubblezone"
)

func (m *M) View() string {
	return zone.Scan(m.metadata.View())
}
