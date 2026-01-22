package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/node/table"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type R source.S

func (r R) Title() string {
	var parts []string
	for i, t := range source.S(r).Titles() {
		if i == 0 || t.Localization() == "en" || t.Localization() == "" {
			parts = append(parts,
				fmt.Sprintf(
					"%v %v",
					lipgloss.NewStyle().Bold(len(parts) == 0).Foreground(map[bool]lipgloss.Color{
						true:  color_profile.ForegroundNormal,
						false: color_profile.ForegroundNegligible,
					}[len(parts) == 0]).Render(t.Title()),
					lipgloss.NewStyle().Foreground(
						color_profile.SupplementaryText,
					).Render(t.Localization()),
				),
			)
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (r R) ID() string {
	id := source.S(r).Header().ID()
	if source.S(r).Header().API() == epb.SourceAPI_SOURCE_API_MAL {
		if source.S(r).Header().Type() == epb.SourceType_SOURCE_TYPE_BOOK_MANGA || source.S(r).Header().Type() == epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL {
			id = fmt.Sprintf("manga/%v", id)
		}
		if source.S(r).Header().Type() == epb.SourceType_SOURCE_TYPE_SERIES_ANIME || source.S(r).Header().Type() == epb.SourceType_SOURCE_TYPE_MOVIE_ANIME {
			id = fmt.Sprintf("anime/%v", id)
		}
	}

	color := map[bool]lipgloss.Color{
		false: color_profile.ForegroundNegligible,
		true:  color_profile.UserViewText,
	}[source.S(r).Header().API() == epb.SourceAPI_SOURCE_API_TRUFFLE || source.S(r).Header().API() == epb.SourceAPI_SOURCE_API_NONE]

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Foreground(color).Render(table.R(r).API()),
		map[bool]string{
			true:  "",
			false: lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(" > "),
		}[source.S(r).Header().ID() == ""],
		lipgloss.NewStyle().Foreground(color).Render(id),
	)
}

func (r R) Notes() string        { return source.S(r).Notes() }
func (r R) Synopsis() string     { return source.S(r).Synopsis() }
func (r R) Genres() string       { return render_tags("Genres", source.S(r).Genres()) }
func (r R) Studios() string      { return render_tags("Studio", source.S(r).Studios()) }
func (r R) Seasons() string      { return render_tags("Seasons", source.S(r).Seasons()) }
func (r R) Authors() string      { return render_tags("Authors", source.S(r).Authors()) }
func (r R) Illustrators() string { return render_tags("Illustrators", source.S(r).Illustrators()) }

func (r R) Status() string {
	return lipgloss.NewStyle().Foreground(color_profile.UserViewText).Render(table.R(r).Status())
}

func (r R) Type() string {
	return lipgloss.NewStyle().Foreground(color_profile.UserViewText).Render(table.R(r).Type())
}

func (r R) Score() string {
	return fmt.Sprintf(
		"%v %v",
		lipgloss.NewStyle().Foreground(color_profile.UserViewText).Render(table.R(r).Score()),
		lipgloss.NewStyle().Foreground(color_profile.ForegroundNegligible).Render(
			fmt.Sprintf("(%d)", source.S(r).Score()),
		),
	)
}

func render_tags(h string, vs []string) string {
	var parts []string
	for _, v := range vs {
		parts = append(parts, lipgloss.NewStyle().Foreground(color_profile.UserViewText).Render(v))
	}
	return strings.Join(
		parts,
		lipgloss.NewStyle().Foreground(
			color_profile.ForegroundNegligible,
		).Render(", "),
	)
}

type RenderHeaderMode int

const (
	RenderHeaderModeNone RenderHeaderMode = iota
	RenderHeaderModeInline
	RenderHeaderModeIndent
	RenderHeaderModeSpacer
)

func with_na(v string) string {
	return map[bool]string{
		false: v,
		true: lipgloss.NewStyle().Foreground(
			color_profile.ForegroundNegligible,
		).Render("N/A"),
	}[lipgloss.Width(v) == 0]
}

func WithHeader(h string, v string, m RenderHeaderMode) string {
	switch m {
	case RenderHeaderModeNone:
		return v
	case RenderHeaderModeInline:
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			fmt.Sprintf(
				"%v: %v",
				lipgloss.NewStyle().Bold(true).Render(h),
				with_na(v),
			),
		)
	case RenderHeaderModeIndent:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render(h),
			lipgloss.NewStyle().Padding(0, 0, 0, 1).Render(with_na(v)),
		)
	case RenderHeaderModeSpacer:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Margin(0, 0, 1, 0).Bold(true).Render(h),
			with_na(v),
		)
	}
	return ""
}
