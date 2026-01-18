package util_node

import (
	"strings"

	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

var (
	Table = table{}
)

type N interface {
	Header() node.H
	Sources() []source.S
	PB() *dpb.Node
	Virtual() (source.S, error)
}

type table struct{}

func (t table) Title(v source.T) string { return v.Title() }

func (t table) API(v epb.SourceAPI) string {
	custom := map[epb.SourceAPI]string{
		epb.SourceAPI_SOURCE_API_NONE: "",
	}
	if s, ok := custom[v]; ok {
		return s
	}
	return strings.ToUpper(
		strings.ReplaceAll(
			strings.TrimPrefix(v.String(), "SOURCE_API_"), "_", " ",
		),
	)
}

func (t table) Type(v epb.SourceType) string {
	custom := map[epb.SourceType]string{
		epb.SourceType_SOURCE_TYPE_SERIES_ANIME:     "Anime",
		epb.SourceType_SOURCE_TYPE_MOVIE_ANIME:      "Anime Movie",
		epb.SourceType_SOURCE_TYPE_BOOK_MANGA:       "Manga",
		epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL: "Light Novel",
	}
	if s, ok := custom[v]; ok {
		return s
	}
	return cases.Title(language.English).String(
		strings.ReplaceAll(
			strings.TrimPrefix(v.String(), "SOURCE_TYPE_"), "_", " ",
		),
	)
}

func (t table) Status(v epb.SourceStatus) string {
	custom := map[epb.SourceStatus]string{
		epb.SourceStatus_SOURCE_STATUS_UNKNOWN: "",
	}
	if s, ok := custom[v]; ok {
		return s
	}
	return cases.Title(language.English).String(
		strings.ReplaceAll(
			strings.TrimPrefix(v.String(), "SOURCE_STATUS_"), "_", " ",
		),
	)
}

func (t table) Score(v int) string {
	score := (v + 10) / 20
	return strings.Repeat("★", score) + strings.Repeat("☆", 5-score)
}
