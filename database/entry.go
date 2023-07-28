package database

import (
	"context"
	"fmt"

	"github.com/minkezhang/truffle/api/graphql/model"
)

type Entry struct {
	data map[string]*model.Entry
}

func NewEntry() *Entry {
	return &Entry{
		data: map[string]*model.Entry{},
	}
}

func (db *Entry) Get(ctx context.Context, id string) (*model.Entry, error) {
	if e, ok := db.data[id]; !ok {
		return nil, fmt.Errorf("cannot find entry: %s", id)
	} else {
		return e, nil
	}
}

func (db *Entry) Put(ctx context.Context, e *model.Entry) (*model.Entry, error) {
	db.data[e.ID] = e
	return e, nil
}

func (db *Entry) Delete(ctx context.Context, id string) (*model.Entry, error) {
	e := db.data[id]
	delete(db.data, id)
	return e, nil
}

func (db *Entry) List(ctx context.Context) (*model.Entry, error) {
	return nil, fmt.Errorf("unimplemented")
}
