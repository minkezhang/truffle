// Package grid defines a set of utility functions to calculate column layouts.
//
// | m | p | content | p | m |

package grid

type C struct {
	Width   int
	Content int
	Margin  int
	Padding int
}

type G struct {
	Width    int
	NColumns int
}

func (g G) Column(n int, m int, p int) C {
	if n == 0 {
		return C{}
	}
	if n > g.NColumns {
		n = g.NColumns
	}
	w := n * g.Width / g.NColumns
	c := w - 2*(m+p)
	return C{
		Width:   w,
		Content: c,
		Margin:  m,
		Padding: p,
	}
}
