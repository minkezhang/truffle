package row

import (
	"github.com/lrstanley/bubblezone"
)

type R struct {
	id        string
	parent_id string
	min       int
	max       int
	width     int
	columns   []*column.C
}

func New(min_width int, max_width int) *R {
	return &R{
		id:  zone.NewPrefix(),
		min: min_width,
		max: max_width,
	}
}

func (r *R) Init() tea.Cmd { return nil }

func (r *R) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return r, nil
}

func (r *R) View() string { return "" }
