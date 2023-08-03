package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"
	"errors"
	"fmt"

	graph "github.com/minkezhang/truffle/api/graphql"
	"github.com/minkezhang/truffle/api/graphql/model"
)

// Sources is the resolver for the sources field.
func (r *metadataResolver) Sources(ctx context.Context, obj *model.Metadata) ([]*model.APIData, error) {
	var sources []*model.APIData
	var errs []error

	for _, s := range obj.Sources {
		// Avoid checking local cache if the data is already populated
		// in e.g. the case of search results.
		//
		// TODO(minkzhang): Change API Get to accept APIData struct.
		if s.Cached {
			sources = append(sources, s)
			continue
		}

		db, ok := r.DB.APIData[s.API]
		if !ok {
			errs = append(errs, fmt.Errorf("unsupported API: %s", s.API))
			continue
		}

		t, err := db.Get(ctx, s.ID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		sources = append(sources, t)
	}
	return sources, errors.Join(errs...)
}

// Metadata returns graph.MetadataResolver implementation.
func (r *Resolver) Metadata() graph.MetadataResolver { return &metadataResolver{r} }

type metadataResolver struct{ *Resolver }
