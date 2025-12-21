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
		Base:    model_ui.New(o.WithNTabs(len(o.Values))),
		sources: sources,
		order:   order,
		kLeft:   zone.NewPrefix(),
		kRight:  zone.NewPrefix(),
	}
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := append([]tea.Cmd{
		m.Base.Update(msg),
	})

	switch msg := msg.(type) {
	case model_ui.FocusMsg:
		if m.ID() == msg.ID {
			cmds = append(cmds, model_ui.ToCommand(SelectMsg{
				Focus: m.sources[m.order[m.Index()]],
			}))
		}
	case tea.KeyMsg:
		if !m.Focus() {
			return m, nil
		}
		switch t := msg.Type; t {
		case tea.KeyLeft:
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{
					ID: m.ID(),
				},
				Index: m.Index() - 1,
			}))
		case tea.KeyRight:
			cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
				BaseMsg: model_ui.BaseMsg{
					ID: m.ID(),
				},
				Index: m.Index() + 1,
			}))
		}
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress {
			if zone.Get(m.kLeft).InBounds(msg) {
				cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
					BaseMsg: model_ui.BaseMsg{
						ID: m.ID(),
					},
					Index: m.Index() - 1,
				}))
			} else if zone.Get(m.kRight).InBounds(msg) {
				cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
					BaseMsg: model_ui.BaseMsg{
						ID: m.ID(),
					},
					Index: m.Index() + 1,
				}))
			} else {
				for i, k := range m.order {
					if zone.Get(k).InBounds(msg) {
						cmds = append(cmds, model_ui.ToCommand(model_ui.FocusMsg{
							BaseMsg: model_ui.BaseMsg{
								ID: m.ID(),
							},
							Index: i,
						}))
					}
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
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
		active := (m.Index() == i)
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
		} else if i > m.Index() {
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
	active := m.Index() == i
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
