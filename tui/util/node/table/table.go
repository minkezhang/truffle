package table

import (
	"strings"

	"github.com/minkezhang/truffle-api/data/source"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type R source.S

func (r R) Title() string { return source.S(r).Title().Title() }

func (r R) API() string {
	custom := map[epb.SourceAPI]string{
		epb.SourceAPI_SOURCE_API_NONE: "",
	}
	if s, ok := custom[source.S(r).Header().API()]; ok {
		return s
	}
	return strings.ToUpper(
		strings.ReplaceAll(
			strings.TrimPrefix(source.S(r).Header().API().String(), "SOURCE_API_"), "_", " ",
		),
	)
}

func (r R) Type() string {
	custom := map[epb.SourceType]string{
		epb.SourceType_SOURCE_TYPE_SERIES_ANIME:     "Anime",
		epb.SourceType_SOURCE_TYPE_MOVIE_ANIME:      "Anime Movie",
		epb.SourceType_SOURCE_TYPE_BOOK_MANGA:       "Manga",
		epb.SourceType_SOURCE_TYPE_BOOK_LIGHT_NOVEL: "Light Novel",
	}
	if s, ok := custom[source.S(r).Header().Type()]; ok {
		return s
	}
	return cases.Title(language.English).String(
		strings.ReplaceAll(
			strings.TrimPrefix(source.S(r).Header().Type().String(), "SOURCE_TYPE_"), "_", " ",
		),
	)
}

func (r R) Status() string {
	custom := map[epb.SourceStatus]string{
		epb.SourceStatus_SOURCE_STATUS_UNKNOWN: "",
	}
	if s, ok := custom[source.S(r).Status()]; ok {
		return s
	}
	return cases.Title(language.English).String(
		strings.ReplaceAll(
			strings.TrimPrefix(source.S(r).Status().String(), "SOURCE_STATUS_"), "_", " ",
		),
	)
}

func (r R) Score() string {
	score := (source.S(r).Score() + 10) / 20
	return strings.Repeat("★", score) + strings.Repeat("☆", 5-score)
}
