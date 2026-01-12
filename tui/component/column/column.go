package column

import (
	"fmt"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/util/render"
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
	id             string
	content        int
	padding        [4]int
	margin         [4]int
	border         lipgloss.Border
	display_border [4]bool
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

func (c *C) SetBorder(b lipgloss.Border, _top, _right, _bottom, _left bool) error {
	content := c.content
	if c.display_border[left] {
		content += c.border.GetLeftSize()
	}
	if c.display_border[right] {
		content += c.border.GetRightSize()
	}
	if _left {
		content -= b.GetLeftSize()
	}
	if _right {
		content -= b.GetRightSize()
	}
	if content < 0 {
		return fmt.Errorf("invalid border (top = %d, right = %d, bottom = %d, left = %d): resizing column results in a negative content block size %d < 0", b.GetTopSize(), b.GetRightSize(), b.GetBottomSize(), b.GetLeftSize(), content)
	}
	c.content = content
	c.border = b
	c.display_border = [4]bool{_top, _right, _bottom, _left}
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
	return lipgloss.NewStyle().Padding(c.padding[:]...).Margin(c.margin[:]...).Border(c.border, c.display_border[:]...).Width(c.Content()).MaxWidth(c.Content())
}

func (c *C) Validate(s string) error {
	err := render.Validate(s, c.Content())
	if err != nil {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			err = fmt.Errorf("%s:%d: %v", file, line, err)
		}
	}
	return err
}

func (c *C) RenderOrDie(s string) string {
	if err := c.Validate(s); err != nil {
		errors.Error(err)
		return ""
	}
	return s
}
