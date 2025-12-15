package title

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
)

var (
	priority = map[string]int{
		"en": 0,
		"":   1,
		"ja": 2,
	}
)

type O struct {
	Titles []atom.T
}

type M struct {
	titles []atom.T
}

func New(o O) *M {
	m := &M{
		titles: append([]atom.T{}, o.Titles...),
	}
	f := func(i, j int) bool {
		if m.titles[i].Localization == m.titles[j].Localization {
			return m.titles[i].Title < m.titles[j].Title
		}
		u, ok := priority[m.titles[i].Localization]
		if !ok {
			u = int(^uint(0) >> 1)
		}
		v, ok := priority[m.titles[j].Localization]
		if !ok {
			v = int(^uint(0) >> 1)
		}
		return u < v
	}
	sort.SliceStable(m.titles, f)
	return m
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return nil, nil }

func (m *M) View() string {
	var parts []string
	for _, t := range m.titles {
		parts = append(parts, strings.Join([]string{
			t.Title,
			lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(t.Localization),
		}, " "))
	}
	return strings.Join(parts, "\n")
}
