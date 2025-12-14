package atom

import (
	"bytes"
	"image"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"

	book_ui "github.com/minkezhang/truffle/tui/component/metadata/book"
)

type O struct {
	Atom *atom.A
}

type M struct {
	atom     *atom.A
	metadata *book_ui.M
	image    string // cache
}

func data(url string) []byte {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	return buf
}

// format converts an image to a string representation of an image.
// From https://github.com/knipferrc/fm.
// TODO(minkezhang): Rewrite.
func format(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.Lanczos)
	b := img.Bounds()
	imageWidth := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for heightCounter := 0; heightCounter < h; heightCounter += 2 {
		for x := imageWidth; x < width; x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}

func New(o O) *M {
	return &M{
		atom: o.Atom,
		metadata: book_ui.New(book_ui.O{
			Book: o.Atom.Metadata().(*book.M),
		}),
	}
}

func (m *M) Init() tea.Cmd {
	return tea.Batch(
		// Load image asynchronously.
		// TOOD(minkezhang): Implement graphics via Kitty when available
		// https://github.com/charmbracelet/bubbletea/issues/163.
		func() tea.Msg {
			// TODO(minkezhang): Cache data and read from cache
			log.Printf("getting url: %v\n", m.atom.PreviewURL())
			d := data("https://cdn.myanimelist.net/images/manga/3/201154l.jpg")
			if d != nil {
				img, _, err := image.Decode(bytes.NewReader(d))
				if err == nil {
					m.image = format(30, img)
				}
				// TODO(minkezhang): Use placeholder image if
				// error
			}
			return ""
		},
		m.metadata.Init(),
	)
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.metadata.Update(msg)

	return m, nil
}

func (m *M) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.image,
			m.metadata.View(),
		),
		lipgloss.NewStyle().Width(80).Render(m.atom.Synopsis()),
	)
}
