package source_list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type V struct {
	API epb.API
	ID  string
}

// Emitted to parents upon focus change.
type SelectMsg struct {
	Blur  V
	Focus V
}

func (v V) String() string {
	switch api := v.API; api {
	case epb.API_API_VIRTUAL:
		return "TRUFFLE"
	default:
		return fmt.Sprintf("%s/%s", strings.ReplaceAll(api.String(), "API_", ""), v.ID)
	}

}

type M struct {
	sources map[string]V // { key: V }
	order   []string     // keys
	index   int
}

type O struct {
	Values []V
}

func New(o O) *M {
	sources := map[string]V{}
	order := []string{}

	for _, v := range o.Values {
		k := zone.NewPrefix()
		sources[k] = v
		order = append(order, k)
	}

	return &M{
		sources: sources,
		order:   order,
		index:   0,
	}
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft, tea.KeyShiftTab:
			src := m.order[m.index]
			m.index = (m.index - 1) % len(m.order)
			if m.index < 0 {
				m.index = len(m.order) - 1
			}
			dst := m.order[m.index]
			return m, func() tea.Msg {
				return SelectMsg{
					Blur:  m.sources[src],
					Focus: m.sources[dst],
				}
			}
		case tea.KeyRight, tea.KeyTab:
			src := m.order[m.index]
			m.index = (m.index + 1) % len(m.order)
			dst := m.order[m.index]
			return m, func() tea.Msg {
				return SelectMsg{
					Blur:  m.sources[src],
					Focus: m.sources[dst],
				}
			}
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft {
			for i, k := range m.order {
				if i != m.index && zone.Get(k).InBounds(msg) {
					src := m.order[m.index]
					dst := m.order[i]
					m.index = i
					return m, func() tea.Msg {
						return SelectMsg{
							Blur:  m.sources[src],
							Focus: m.sources[dst],
						}
					}
				}
			}
		}
	}

	return m, nil
}

var (
	// { active: lipgloss.Style }
	borders = map[bool]lipgloss.Style{
		false: lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Border(
			lipgloss.NormalBorder(), false, false, true, false,
		).BorderForeground(lipgloss.Color("8")).Foreground(lipgloss.Color("8")),
		true: lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Border(
			lipgloss.NormalBorder(), false, false, true, false,
		).BorderForeground(lipgloss.Color("6")).Foreground(lipgloss.Color("6")),
	}
)

func (m *M) View() string {
	parts := []string{}
	for i, k := range m.order {
		active := (m.index == i)
		parts = append(
			parts,
			zone.Mark(k, borders[active].Render(m.sources[k].String())),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		parts...,
	)
}
