package grid

import (
	"testing"
)

func TestWithPadding(t *testing.T) {
	c := C{
		Content: 100,
		Margin:  1,
		Padding: 1,
	}
	want := C{
		Content: 98,
		Margin:  1,
		Padding: 2,
	}
	if got := c.WithPadding(2); got != want {
		t.Errorf("WithPadding() = %v, want = %v", got, want)
	}
}

func TestColumn(t *testing.T) {
	configs := []struct {
		name string
		g    G
		n    int
		m    int
		p    int
		want C
	}{
		{
			name: "C=2/N=1",
			g: G{
				Width: 100,
				N:     2,
			},
			n: 1,
			m: 0,
			p: 0,
			want: C{
				Content: 50,
				Margin:  0,
				Padding: 0,
			},
		},
		{
			name: "C=3/N=2",
			g: G{
				Width: 300,
				N:     3,
			},
			n: 2,
			m: 0,
			p: 0,
			want: C{
				Content: 200,
				Margin:  0,
				Padding: 0,
			},
		},
		{
			name: "C=3/N=3",
			g: G{
				Width: 300,
				N:     3,
			},
			n: 3,
			m: 0,
			p: 0,
			want: C{
				Content: 300,
				Margin:  0,
				Padding: 0,
			},
		},
		{
			name: "C=3/N=2/M/P",
			g: G{
				Width: 300,
				N:     3,
			},
			n: 2,
			m: 1,
			p: 1,
			want: C{
				Content: 196,
				Margin:  1,
				Padding: 1,
			},
		},
	}
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.g.Column(c.n, c.m, c.p); got != c.want {
				t.Errorf("Column() = %v, want = %v", got, c.want)
			}
		})
	}
}
