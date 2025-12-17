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

// L defines the column layout
type L struct {
	Content int
	Margin  int
	Padding int
}

func (l L) Width() int { return l.Content + 2*(l.Margin+l.Padding) }

func (l L) Style() lipgloss.Style {
	return lipgloss.NewStyle().Width(
		l.Content+2*l.Padding,
	).Margin(0, l.Margin).Padding(0, l.Padding)
}

func (l L) WithPadding(n int) L {
	return L{
		Content: l.Content - 2*(n-l.Padding),
		Margin:  l.Margin,
		Padding: n,
	}
}

func (l L) WithMargin(n int) L {
	return L{
		Content: l.Content - 2*(n-l.Margin),
		Margin:  n,
		Padding: l.Padding,
	}
}

func (l L) Grid(n int) G {
	return G{
		Width:    l.Content,
		NColumns: n,
	}
}

type G struct {
	Width    int
	NColumns int
}

func (g G) Column(n int, m int, p int) L {
	if n == 0 {
		return L{}
	}
	if n > g.NColumns {
		n = g.NColumns
	}
	w := n * g.Width / g.NColumns
	c := w - 2*(m+p)
	return L{
		Content: c,
		Margin:  m,
		Padding: p,
	}
}
