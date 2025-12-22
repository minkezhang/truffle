// Package grid defines a set of utility functions to calculate column layouts.
//
// | m | p | content | p | m |
//
// lipgloss.Style.Width includes padding and content, but excludes the margins,
// as is the case with CSS.
package grid

import (
	"github.com/charmbracelet/lipgloss"
)

// C defines the column layout
type C struct {
	Content int
	Margin  int
	Padding int
}

func (c C) Width() int { return c.Content + 2*(c.Margin+c.Padding) }

func (c C) Style() lipgloss.Style {
	return lipgloss.NewStyle().Width(
		c.Content,
	).Margin(0, c.Margin).Padding(0, c.Padding)
}

func (c C) WithPadding(n int) C {
	return C{
		Content: c.Content - 2*(n-c.Padding),
		Margin:  c.Margin,
		Padding: n,
	}
}

func (c C) WithMargin(n int) C {
	return C{
		Content: c.Content - 2*(n-c.Margin),
		Margin:  n,
		Padding: c.Padding,
	}
}

func (c C) Grid(n int) G {
	return G{
		Width: c.Content,
		N:     n,
	}
}

type G struct {
	Width int
	N     int
}

func (g G) Column(n int) C {
	if n == 0 {
		return C{}
	}
	if n > g.N {
		n = g.N
	}
	return C{
		Content: n * g.Width / g.N,
	}
}

func GetFrame(s lipgloss.Style) int {
	return s.GetBorderLeftSize() + s.GetBorderRightSize() + s.GetPaddingLeft() + s.GetPaddingRight() + s.GetMarginLeft() + s.GetMarginRight()
}
