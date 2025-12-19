package title

import (
	"fmt"
	"reflect"
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

type updateTitlesMsg struct {
	model_ui.BaseMsg
	payload []atom.T
}

func UpdateTitlesAsync(m tea.Model, v []atom.T) tea.Cmd {
	if m, ok := m.(*M); ok {
		return func() tea.Msg {
			f := func(i, j int) bool {
				if v[i].Localization == v[j].Localization {
					return v[i].Title < v[j].Title
				}
				u, ok := priority[v[i].Localization]
				if !ok {
					u = int(^uint(0) >> 1)
				}
				w, ok := priority[v[j].Localization]
				if !ok {
					w = int(^uint(0) >> 1)
				}
				return u < w
			}
			v = append([]atom.T{}, v...)
			sort.SliceStable(v, f)
			return updateTitlesMsg{
				BaseMsg: model_ui.BaseMsg{
					ID: m.ID(),
				},
				payload: v,
			}
		}
	}

	return model_ui.ErrorCmd(fmt.Errorf("incorrect model type: %v", reflect.TypeOf(m)))
}

func New(o O) *M {
	m := &M{
		titles: append([]atom.T{}, o.Titles...),
		Base:   model_ui.New(o.O),
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

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateTitlesMsg:
		if m.ID() == msg.ID {
			m.titles = msg.payload
		}
	}
	return m, nil
}

func (m *M) View() string {
	var parts []string
	for i, t := range m.titles[:2] {
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
