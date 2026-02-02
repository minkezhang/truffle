package focusable

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

type Node interface {
	base.Identifiable

	// OnFocus may be called in the Update method by the directory. This
	// function directly instructs the node to activate any sub-elements for
	// input, e.g. focusing a specific input field.
	OnFocus(i int) tea.Cmd

	// OnBlur may be called in the Update method by the directory. This
	// function directly instructs the node to deactivate all sub-elements.
	OnBlur() tea.Cmd

	// NElements returns the number of interactable (and therefore
	// tab-focusable) elements in this node. NElements does not take into
	// account the number of [Node] children.
	NElements() int

	// IsInvisible indicates the element is not currently rendered, i.e. is
	// hidden behind a collapsed zippy. This is different from an element
	// being out of frame (e.g. in a viewport).
	IsInvisible() bool
	SetIsInvisible(v bool) tea.Cmd

	FocusIndex() int
	FocusState() types.FocusState
}

type D struct {
	directory *base.D

	current_node_id string
	dirty           bool
	order_cache     []string
}

func New() *D {
	return &D{
		directory: base.New(),
		dirty:     true,
	}
}

// order returns the pre-order traversal of the directory.
func (d *D) order() []string {
	if !d.dirty {
		return d.order_cache
	}
	var f func(n string) []string
	f = func(n string) []string {
		var result = []string{}
		if n != "" {
			result = append(result, n)
		}
		for _, c := range d.directory.Children[n] {
			result = append(result, f(c)...)
		}
		return result
	}
	d.order_cache = f("")
	return d.order_cache
}

func (d *D) Init() tea.Cmd { return nil }

func (d *D) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if _, ok := d.directory.Nodes[d.current_node_id]; ok {
			switch t := msg.Type; t {
			case tea.KeyShiftTab:
				dst := d.directory.Nodes[d.current_node_id].(Node).FocusIndex() - 1
				cmds = append(cmds, d.directory.Nodes[d.current_node_id].(Node).OnFocus(dst))
			case tea.KeyTab:
				dst := d.directory.Nodes[d.current_node_id].(Node).FocusIndex() + 1
				cmds = append(cmds, d.directory.Nodes[d.current_node_id].(Node).OnFocus(dst))
			}
		}
	case base.RegisterMessage:
		if n, ok := msg.Node.(Node); ok {
			_, c := d.directory.Update(msg)
			cmds = append(cmds, c)

			d.dirty = true

			if d.current_node_id == "" && !n.IsInvisible() {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
		}
	case base.PutChildrenMessage:
		d.dirty = true
		_, c := d.directory.Update(msg)
		cmds = append(cmds, c)

	case types.FocusMessage:
		if n, ok := d.directory.Nodes[d.current_node_id]; ok {
			cmds = append(cmds, n.(Node).OnBlur())
		}
		if n, ok := d.directory.Nodes[msg.ID]; ok {
			cmds = append(cmds, n.(Node).OnFocus(msg.Index))
		}
		d.current_node_id = msg.ID
	case types.EOFMessage:
		target_index := slices.IndexFunc(d.order(), func(v string) bool { return v == d.current_node_id })
		var target_id string

		// Skip any "invisible" nodes which are hidden away beneath
		// zippy containers.
		if msg.IsHead {
			for target_index := target_index - 1; ; target_index -= 1 {
				if target_index < 0 {
					target_index = len(d.order()) - 1
				}
				target_id = d.order()[target_index]
				if !d.directory.Nodes[target_id].(Node).IsInvisible() {
					break
				}
			}
		} else {
			for target_index := target_index + 1; ; target_index += 1 {
				if target_index >= len(d.order()) {
					target_index = 0
				}
				target_id = d.order()[target_index]
				if !d.directory.Nodes[target_id].(Node).IsInvisible() {
					break
				}
			}
		}

		focus_index := 0
		if msg.IsHead {
			focus_index = d.directory.Nodes[target_id].(Node).NElements() - 1
		} else {
			focus_index = 0
		}
		if _, ok := d.directory.Nodes[d.current_node_id]; ok {
			cmds = append(cmds, d.directory.Nodes[d.current_node_id].(Node).OnBlur())
		}
		if _, ok := d.directory.Nodes[target_id]; ok {
			cmds = append(cmds, d.directory.Nodes[target_id].(Node).OnFocus(focus_index))
		}
		d.current_node_id = target_id
	}
	return d, tea.Batch(cmds...)
}

func (d *D) View() string {
	var tree func(indent int, prefix string, nodes []string) []string
	// From github.com/campoy/tools/tree.
	tree = func(indent int, prefix string, nodes []string) []string {
		result := []string{}
		for i, n := range nodes {
			color := color_profile.ForegroundNormal

			node, ok := d.directory.Nodes[n]
			if ok {
				if node.(Node).FocusState() == types.FocusStateActive {
					color = color_profile.ForegroundCritical
				} else if !node.(Node).IsInvisible() {
					color = color_profile.ForegroundImportant
				} else if node.(Node).IsInvisible() && d.directory.Children[n] == nil {
					color = color_profile.ForegroundNegligible
				}
			}
			directory := "│   "
			file := "├── "
			if i == len(nodes)-1 {
				directory = "    "
				file = "└── "
			}
			if n == "" {
				directory = " "
				result = append(result, "(root)")
			} else {
				result = append(
					result,
					fmt.Sprintf(
						"%v%v%v",
						lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(prefix),
						lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(file),
						lipgloss.NewStyle().Foreground(color).Render(n),
					),
				)
			}
			result = append(
				result,
				tree(indent+1, prefix+directory, d.directory.Children[n])...,
			)
		}
		return result
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Node Directory"),
			lipgloss.NewStyle().Margin(0, 0, 0, 1).Foreground(color_profile.ForegroundCritical).Render("Active"),
			lipgloss.NewStyle().Margin(0, 0, 0, 1).Foreground(color_profile.ForegroundImportant).Render("Focusable"),
			lipgloss.NewStyle().Margin(0, 0, 0, 1).Foreground(color_profile.ForegroundNormal).Render("No Focusable Elements"),
			lipgloss.NewStyle().Margin(0, 0, 1, 1).Foreground(color_profile.ForegroundNegligible).Render("Leaf"),
		),
		strings.Join(tree(0, "", []string{""}), "\n"),
	)
}
