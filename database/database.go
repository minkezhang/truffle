package database

import (
	"context"

	"github.com/minkezhang/truffle/api/go/database/utils"
	"github.com/minkezhang/truffle/client/mal"
	"github.com/minkezhang/truffle/client/truffle"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	cpb "github.com/minkezhang/truffle/api/go/config"
	dpb "github.com/minkezhang/truffle/api/go/database"
)

type DB struct {
	truffle *truffle.C

	config *cpb.Config
}

func New(db *truffle.C, config *cpb.Config) *DB {
	return &DB{
		truffle: db,
		config:  config,
	}
}

func (db *DB) Add(ctx context.Context, epb *dpb.Entry) (*dpb.Entry, error) {
	return db.truffle.Add(ctx, epb)
}

type GetOpts struct {
	APIs []dpb.API
}

func (db *DB) Get(ctx context.Context, id *dpb.LinkedID, opts interface{}) (*dpb.Entry, error) {
	query, ok := opts.(GetOpts)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid search opts provided")
	}

	apis := map[dpb.API]bool{}
	for _, api := range query.APIs {
		apis[api] = true
	}

	closed := map[string]bool{}
	open := []*dpb.LinkedID{id}

	epb := &dpb.Entry{}

	for len(open) > 0 {
		id, open = open[0], open[1:]
		if closed[utils.ID(id)] {
			continue
		}
		closed[utils.ID(id)] = true

		if f := map[dpb.API]func(context.Context, *dpb.LinkedID, interface{}) (*dpb.Entry, error){
			dpb.API_API_TRUFFLE: db.truffle.Get,
			dpb.API_API_MAL:     mal.New(db.config.GetMal()).Get,
		}[id.GetApi()]; f != nil {
			if !apis[id.GetApi()] {
				continue
			}

			fpb, err := f(ctx, id, nil)
			if err != nil {
				// We only consider user DB lookup errors to be
				// fatal.
				if id.GetApi() == dpb.API_API_TRUFFLE {
					return nil, err
				} else {
					return epb, nil
				}
			}

			// Do not override already populated important fields.
			if epb.GetId() != nil {
				fpb.Id = nil
			}
			if epb.GetScore() != 0 {
				fpb.Score = 0
			}
			if epb.GetEtag() != nil {
				fpb.Etag = nil
			}

			proto.Merge(epb, fpb)

			open = append(open, fpb.GetLinkedIds()...)
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "unsupported API specified for Get()")
		}

	}

	return utils.Clean(epb), nil
}

func (db *DB) Put(ctx context.Context, epb *dpb.Entry) (*dpb.Entry, error) {
	return db.truffle.Put(ctx, epb)
}

func (db *DB) Delete(ctx context.Context, id *dpb.LinkedID) (*dpb.Entry, error) {
	return db.truffle.Delete(ctx, id)
}

type SearchOpts struct {
	Title  string
	Corpus dpb.Corpus

	APIs []dpb.API

	// MAL contains MAL-specific API options.
	MAL mal.SearchOpts
}

func (db *DB) Search(ctx context.Context, opts interface{}) ([]*dpb.Entry, error) {
	query, ok := opts.(SearchOpts)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid search opts provided")
	}

	apis := map[dpb.API]bool{}
	for _, api := range query.APIs {
		apis[api] = true
	}

	var candidates []*dpb.Entry

	if apis[dpb.API_API_TRUFFLE] {
		if cs, err := db.truffle.Search(
			ctx,
			truffle.SearchOpts{
				Title:  query.Title,
				Corpus: query.Corpus,
			},
		); err != nil {
			return nil, err
		} else {
			candidates = append(candidates, cs...)
		}
	}

	duplicates := map[string]bool{}
	for _, epb := range candidates {
		for _, id := range epb.GetLinkedIds() {
			duplicates[utils.ID(id)] = true
		}
	}

	if apis[dpb.API_API_MAL] {
		if cs, err := mal.New(db.config.GetMal()).Search(
			ctx,
			mal.SearchOpts{
				Title:  query.Title,
				Corpus: query.Corpus,
			},
		); err != nil {
			// Only return fatal if MAL is the only API being
			// used; otherwise continue.
			if len(apis) == 1 {
				return nil, err
			}
		} else {
			candidates = append(candidates, cs...)
		}
	}

	var unique []*dpb.Entry
	for _, epb := range candidates {
		// Skip reporting entries which are already accounted for in the
		// user DB (if we are returning user DB results).
		if apis[dpb.API_API_TRUFFLE] && duplicates[utils.ID(epb.GetId())] {
			continue
		}
		unique = append(unique, epb)
	}
	return unique, nil
}

func Marshal(db *DB) ([]byte, error) { return truffle.Marshal(db.truffle) }
