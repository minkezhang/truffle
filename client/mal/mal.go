package mal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/minkezhang/truffle/api/graphql/model"
	"github.com/nstratos/go-myanimelist/mal"
)

type MAL struct {
	manga *Manga
	anime *Anime
}

func New(o O) *MAL {
	c := mal.NewClient(
		&http.Client{
			Transport: &transport{
				ClientID: o.ClientID,
			},
		},
	)
	return &MAL{
		manga: NewManga(c),
		anime: NewAnime(c),
	}
}

func (c *MAL) API() model.APIType { return model.APITypeAPIMal }

func (c *MAL) Get(ctx context.Context, id string) (*model.APIData, error) {
	parts := strings.Split(id, "/")
	corpus := parts[0]
	if corpus == "manga" {
		return c.manga.Get(ctx, parts[1])
	}
	if corpus == "anime" {
		return c.anime.Get(ctx, parts[1])
	}
	return nil, fmt.Errorf("unimplemented MAL corpus: %s", corpus)
}