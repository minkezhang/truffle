package model

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/util/grid"
)

var (
	debug = lipgloss.NewStyle().Background(
		lipgloss.Color("#FF0000")).BorderBackground(
		lipgloss.Color("#00FF00")).MarginBackground(
		lipgloss.Color("#0000FF"))
)

func TestCheck(t *testing.T) {
	configs := []struct {
		name    string
		s       string
		w       int
		success bool
	}{
		{
			name: "All",
			s: debug.Margin(
				0, 1,
			).Padding(
				0, 1,
			).Border(
				lipgloss.NormalBorder(),
			).BorderTop(false).BorderBottom(false).Render("ABC"),
			w:       9,
			success: true,
		},
		{
			name:    "Padding",
			s:       debug.Padding(0, 1).Render("ABC"),
			w:       5,
			success: true,
		},
		{
			name:    "Margin",
			s:       debug.Margin(0, 1).Render("ABC"),
			w:       5,
			success: true,
		},
	}
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			err := check(c.s, grid.C{
				Content: c.w,
			})
			if c.success && err != nil {
				t.Errorf("check() raised unexpected error: %v", err)
			}
			if !c.success && err == nil {
				t.Errorf("check() unexpectedly succeeded")
			}
		})
	}
}
