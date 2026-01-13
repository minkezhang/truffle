package column

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestWidth(t *testing.T) {
	configs := []struct {
		name string
		s    string
		want int
	}{
		{
			name: "Trivial",
			s:    "",
			want: 0,
		},
		{
			name: "Simple",
			s:    "foo",
			want: 3,
		},
		{
			name: "WithPadding",
			s:    lipgloss.NewStyle().Padding(0, 1, 0, 1).Render("foo"),
			want: 5,
		},
		{
			name: "WithMargin",
			s:    lipgloss.NewStyle().Margin(0, 1, 0, 1).Render("foo"),
			want: 5,
		},
		{
			name: "WithBorder",
			s:    lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, true).Render("foo"),
			want: 5,
		},
		{
			name: "ExplicitWidth/WithPadding",
			s:    lipgloss.NewStyle().Padding(0, 1, 0, 1).Width(5).Render("foo"), // content + padding
			want: 5,
		},
		{
			name: "ExplicitWidth/WithMargin",
			s:    lipgloss.NewStyle().Margin(0, 1, 0, 1).Width(3).Render("foo"),
			want: 5,
		},
		{
			name: "ExplicitWidth/WithBorder",
			s:    lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, true, false, true).Width(3).Render("foo"),
			want: 5,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := lipgloss.Width(c.s); got != c.want {
				t.Errorf("Width(%s) = %v, want = %v", c.s, got, c.want)
			}
		})
	}
}
