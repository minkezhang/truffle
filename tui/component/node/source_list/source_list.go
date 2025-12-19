// Package source_list provides a tab bar for a list of API:IDs.
//
// TODO(minkezhang): Implement tests for rendering narrow regions.
// TODO(minkezhang): Extrapolate model styling into a struct to guarantee
// content falls within the width.
// TODO(minkezhang): When at end of list and last element still doesn't fit,
// truncate left-most element instead.
package source_list

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
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
	*model_ui.Base

	sources map[string]V // { key: V }
	order   []string     // keys
	index   int

	kLeft  string
	kRight string
	iStart int
}

type O struct {
	model_ui.O
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
		Base:    model_ui.New(o.O),
		sources: sources,
		order:   order,
		index:   0,
		kLeft:   zone.NewPrefix(),
		kRight:  zone.NewPrefix(),
	}
}

type updateSourcesMsg struct {
	model_ui.BaseMsg
	payload []V
}

func UpdateSourcesAsync(m tea.Model, v []V) tea.Cmd {
	if m, ok := m.(*M); ok {
		return func() tea.Msg {
			return updateSourcesMsg{
				BaseMsg: model_ui.BaseMsg{
					ID: m.ID(),
				},
				payload: v,
			}
		}
	}
	return model_ui.ErrorCmd(fmt.Errorf("incorrect model type: %v", reflect.TypeOf(m)))
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateSourcesMsg:
		if m.ID() == msg.ID {
			sources := map[string]V{}
			order := []string{}
			for _, v := range msg.payload {
				k := zone.NewPrefix()
				sources[k] = v
				order = append(order, k)
				m.sources = sources
				m.order = order
			}
			m.index = 0
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft, tea.KeyShiftTab:
			src := m.order[m.index]
			m.index = (m.index - 1) % len(m.order)
			if m.index < 0 {
				m.index = len(m.order) - 1
			}
			dst := m.order[m.index]
			if m.index < m.iStart {
				m.iStart = m.index
			}
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
			if m.index < m.iStart {
				m.iStart = m.index
			}
			return m, func() tea.Msg {
				return SelectMsg{
					Blur:  m.sources[src],
					Focus: m.sources[dst],
				}
			}
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress {
			if zone.Get(m.kLeft).InBounds(msg) {
				src := m.order[m.index]
				m.index = (m.index - 1) % len(m.order)
				if m.index < 0 {
					m.index = len(m.order) - 1
				}
				dst := m.order[m.index]
				if m.index < m.iStart {
					m.iStart = m.index
				}
				return m, func() tea.Msg {
					return SelectMsg{
						Blur:  m.sources[src],
						Focus: m.sources[dst],
					}
				}
			} else if zone.Get(m.kRight).InBounds(msg) {
				src := m.order[m.index]
				m.index = (m.index + 1) % len(m.order)
				dst := m.order[m.index]
				if m.index < m.iStart {
					m.iStart = m.index
				}
				return m, func() tea.Msg {
					return SelectMsg{
						Blur:  m.sources[src],
						Focus: m.sources[dst],
					}
				}
			} else {
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
	arrows = lipgloss.NewStyle().Padding(0, 0).Margin(0, 1, 1, 1)
)

func (m *M) View() string {
	parts := []string{}
	lengths := []int{}

	for i, k := range m.order {
		active := (m.index == i)
		p := zone.Mark(k, borders[active].Render(m.sources[k].String()))

		parts = append(parts, p)
		lengths = append(lengths, lipgloss.Width(p))
	}

	left := zone.Mark(m.kLeft, arrows.Render("<"))
	right := zone.Mark(m.kRight, arrows.Render(">"))
	ll := lipgloss.Width(left)
	lr := lipgloss.Width(right)

	iStart := m.iStart
	iEnd := m.iStart

	length := ll + lr

	for i, l := range lengths[iStart:] {
		if length+l < m.Column().Content {
			length += l
			iEnd = i
		} else if i == len(lengths)-1 {
			// Won't need the right arrow, discount it from length
			// calculations
			length -= lr
			if length+l < m.Column().Content { // Last tab now wholy fits
				length += l
				iEnd = i
			} else { // Append partial
				iEnd = i
				parts[i] = m.partial(i, m.Column().Content-length)
			}
		} else if i > m.index {
			// The active element is already in the tab list, so add the
			// last partial tab and return
			iEnd = i
			parts[i] = m.partial(i, m.Column().Content-length)
			break
		} else { // Need to advance iStart
			length += (-lengths[iStart] + lengths[i])
			if iStart == 0 {
				length += lipgloss.Width(left)
			}
			iStart += 1
			iEnd = i
			if iEnd == len(lengths)-1 {
				length -= lipgloss.Width(right)
			}
		}
	}

	parts = parts[iStart : iEnd+1]
	if iStart > 0 {
		parts = append([]string{left}, parts...)
	}
	if iEnd < len(lengths)-1 {
		parts = append(parts, right)
	}
	return m.RenderOrDie(
		lipgloss.NewStyle().Width(m.Column().Content).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				parts...,
			),
		),
	)
}

// partial renders a partial tab, truncated on the right.
func (m *M) partial(i int, w int) string {
	active := m.index == i
	k := m.order[i]
	s := m.sources[k].String()
	if w == 0 {
		style := borders[active].Copy().MarginRight(0).PaddingRight(0).PaddingLeft(0).MarginLeft(0)
		return zone.Mark(k, style.Render(""))
	} else if w == 1 {
		style := borders[active].Copy().MarginRight(0).PaddingRight(0).PaddingLeft(0)
		return zone.Mark(k, style.Render(""))
	} else if w < len(s)+3 {
		style := borders[active].Copy().MarginRight(0).PaddingRight(0)
		return zone.Mark(k, style.Render(s[:w-2]))
	} else if w < len(s)+4 {
		style := borders[active].Copy().MarginRight(0)
		return zone.Mark(k, style.Render(s))
	} else {
		return zone.Mark(k, borders[active].Render(s))
	}
}
