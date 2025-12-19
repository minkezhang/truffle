package title

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

var (
	priority = map[string]int{
		"en": 0,
		"":   1,
		"ja": 2,
	}
)

type O struct {
	model_ui.O

	Titles []atom.T
}

type M struct {
	*model_ui.Base

	titles []atom.T
}

func Make(o O) M {
	titles := append([]atom.T{}, o.Titles...)
	f := func(i, j int) bool {
		if titles[i].Localization == titles[j].Localization {
			return titles[i].Title < titles[j].Title
		}
		u, ok := priority[titles[i].Localization]
		if !ok {
			u = int(^uint(0) >> 1)
		}
		v, ok := priority[titles[j].Localization]
		if !ok {
			v = int(^uint(0) >> 1)
		}
		return u < v
	}
	sort.SliceStable(titles, f)

	m := M{
		titles: titles,
		Base:   model_ui.New(o.O),
	}
	return m
}

func (m M) Init() tea.Cmd                           { return nil }
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m M) View() string {
	var parts []string
	titles := m.titles
	if len(m.titles) > 2 {
		titles = titles[:2]
	}
	for i, t := range titles {
		styleTitle := lipgloss.NewStyle()
		if i == 0 {
			styleTitle = styleTitle.Bold(true)
		} else {
			styleTitle = styleTitle.Foreground(lipgloss.Color("8"))
		}
		parts = append(parts, styleTitle.Render(t.Title))
	}
	return m.RenderOrDie(strings.Join(parts, "\n"))
}
