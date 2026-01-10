package renderable

import (
	"github.com/charmbracelet/lipgloss"
)

type WidthOption interface {
	is_option()
}

type FixedWidth int

func (o FixedWidth) is_option() {}

type RelativeWidth float64

func (o RelativeWidth) is_option() {}

func New(min_width int, max_width int, width_options ...WidthOption) *Row {
	return nil
}

type Row struct {
	min_width int
	max_width int
	columns   []*Column
}

type r Rows
