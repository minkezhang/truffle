package column

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
)

type direction int

const (
	top direction = iota // CSS order
	right
	bottom
	left
)

func New(content int) *C {
	return &C{
		id:      zone.NewPrefix(),
		content: content,
	}
}

type C struct {
	id      string
	content int
	padding [4]int
	margin  [4]int
	border  lipgloss.Border
}

func (c *C) ID() string   { return c.id }
func (c *C) Content() int { return c.content }

func (c *C) Padding() (int, int, int, int) {
	return c.padding[top], c.padding[right], c.padding[bottom], c.padding[left]
}

func (c *C) Margin() (int, int, int, int) {
	return c.margin[top], c.margin[right], c.margin[bottom], c.margin[left]
}

func (c *C) Border() (int, int, int, int) {
	return c.border.GetTopSize(), c.border.GetRightSize(), c.border.GetBottomSize(), c.border.GetLeftSize()
}

func (c *C) Width() int {
	return c.padding[left] + c.padding[right] + c.margin[left] + c.margin[right] + c.border.GetLeftSize() + c.border.GetRightSize()
}

func (c *C) SetWidth(w int) error {
	content := w - c.padding[left] - c.padding[right] - c.border.GetLeftSize() - c.border.GetRightSize() - c.margin[left] - c.margin[right]
	if content < 0 {
		return fmt.Errorf("invalid width %d: resizing column results in a negative content block size %d < 0", w, content)
	}
	c.content = content
	return nil
}

func (c *C) SetPadding(top, right, bottom, left int) error {
	content := c.content + c.padding[left] + c.padding[right] - left - right
	if content < 0 {
		return fmt.Errorf("invalid padding (top = %d, right = %d, bottom = %d, left = %d): resizing column results in a negative content block size %d < 0", top, right, bottom, left, content)
	}
	c.content = content
	c.padding = [4]int{top, right, bottom, left}
	return nil
}

func (c *C) SetBorder(b lipgloss.Border) error {
	content := c.content + c.border.GetLeftSize() + c.border.GetRightSize() - b.GetLeftSize() - b.GetRightSize()
	if content < 0 {
		return fmt.Errorf("invalid border (top = %d, right = %d, bottom = %d, left = %d): resizing column results in a negative content block size %d < 0", b.GetTopSize(), b.GetRightSize(), b.GetBottomSize(), b.GetLeftSize(), content)
	}
	c.content = content
	c.border = b
	return nil
}

func (c *C) SetMargin(top, right, bottom, left int) error {
	content := c.content + c.margin[left] + c.margin[right] - left - right
	if content < 0 {
		return fmt.Errorf("invalid margin (top = %d, right = %d, bottom = %d, left = %d): resizing column results in a negative content block size %d < 0", top, right, bottom, left, content)
	}
	c.content = content
	c.margin = [4]int{top, right, bottom, left}
	return nil
}

func (c *C) Style() lipgloss.Style {
	return lipgloss.NewStyle().Padding(
		c.padding[top], c.padding[right], c.padding[bottom], c.padding[left],
	).Margin(
		c.margin[top], c.margin[right], c.margin[bottom], c.padding[left],
	).Border(c.border)
}

func with_xy_axis(s string) string {
	w := lipgloss.Width(s)
	var axis string
	if w < 10 {
		axis = []string{
			"",
			"│",
			"├─",
			"├──",
			"├───",
			"├────",
			"├────┐",
			"├────┬─",
			"├────┬──",
			"├────┬───",
		}[w]
	} else {
		axis = fmt.Sprintf(
			"%v%v%v",
			"├────┬────",
			strings.Repeat(
				"┼────┬────",
				(w/10)-1,
			),
			[]string{
				"",
				"┤",
				"┼─",
				"┼──",
				"┼───",
				"┼────",
				"┼────┐",
				"┼────┬─",
				"┼────┬──",
				"┼────┬───",
			}[w%10],
		)
	}
	return fmt.Sprintf("  %v\n%v", axis, with_y_axis(s))

}

func with_y_axis(s string) string {
	h := lipgloss.Height(s)
	var axis string
	if h < 10 {
		axis = []string{
			"",
			"─",
			"┬│",
			"┬││",
			"┬│││",
			"┬││││",
			"┬││││└",
			"┬││││├│",
			"┬││││├││",
			"┬││││├│││",
		}[h]
	} else {
		axis = fmt.Sprintf(
			"%v%v%v",
			"┬││││├││││",
			strings.Repeat(
				"┼││││├││││",
				(h/10)-1,
			),
			[]string{
				"",
				"┴",
				"┼│",
				"┼││",
				"┼│││",
				"┼││││",
				"┼││││└",
				"┼││││├│",
				"┼││││├││",
				"┼││││├│││",
			}[h%10],
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		strings.Join(strings.Split(axis, ""), "\n"),
		" ",
		s,
	)
}

func check(s string, max_width int) error {
	w := lipgloss.Width(s)
	if w <= max_width {
		return nil
	}
	return fmt.Errorf(
		"rendered string exceeded bounding box: %d > %d\n  %s\n%s\n",
		w,
		max_width,
		strings.Repeat(" ", max_width)+"↓",
		with_xy_axis(s),
	)
}

func (c *C) Check(s string) error {
	err := check(s, c.Content())
	if err != nil {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			err = fmt.Errorf("%s:%d: %v", file, line, err)
		}
	}
	return err
}
