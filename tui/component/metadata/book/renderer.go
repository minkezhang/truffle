package book

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"
)

var (
	magenta = lipgloss.NewStyle().Foreground(
		lipgloss.Color("5"),
	)
)

type R struct{}

func (r R) BookType(m *book.M) string {
	return map[bool]map[bool]string{
		true: map[bool]string{
			true:  "manga",
			false: "light novel",
		},
		false: map[bool]string{
			true:  "comic",
			false: "book",
		},
	}[m.IsManga()][m.IsIllustrated()]
}

func (r R) Genres(m *book.M) []string {
	res := []string{}
	for _, g := range m.Genres() {
		res = append(res, magenta.Render(g))
	}
	return res
}

func (r R) Authors(m *book.M) []string {
	res := []string{}
	for _, a := range m.Authors() {
		res = append(res, magenta.Render(
			strings.TrimSpace(a)))
	}
	return res
}

func (r R) Illustrators(m *book.M) []string {
	res := []string{}
	for _, i := range m.Illustrators() {
		res = append(res, magenta.Render(
			strings.TrimSpace(i)))
	}
	return res
}

func (r R) LastUpdated(m *book.M) string {
	return magenta.Render(m.LastUpdated().Format("2006-01-02"))
}
