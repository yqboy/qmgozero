package qmg

import (
	"context"
	"log"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

type Model struct {
	db   Database
	opts []options.ClientOptions
}

func MustNewModel(url, database string, opts ...options.ClientOptions) *Model {
	model, err := NewModel(url, database, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return model
}

func NewModel(url, database string, opts ...options.ClientOptions) (*Model, error) {
	ctx := context.Background()
	cli, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: url}, opts...)
	if err != nil {
		return nil, err
	}

	db := cli.Database(database)

	return &Model{
		db:   newCli(ctx, cli, db),
		opts: opts,
	}, nil
}

func (m *Model) Insert(coll string, doc interface{}, opts ...options.InsertOneOptions) (string, error) {
	return m.db.Insert(coll, doc, opts...)
}

func (m *Model) InsertMany(coll string, doc interface{}, opts ...options.InsertManyOptions) ([]string, error) {
	return m.db.InsertMany(coll, doc, opts...)
}

func (m *Model) Remove(coll string, selector interface{}, opts ...options.RemoveOptions) error {
	return m.db.Remove(coll, selector, opts...)
}
func (m *Model) RemoveAll(coll string, selector interface{}, opts ...options.RemoveOptions) error {
	return m.db.RemoveAll(coll, selector, opts...)
}

func (m *Model) Update(coll string, selector, update interface{}, opts ...options.UpdateOptions) error {
	return m.db.Update(coll, selector, update, opts...)
}

func (m *Model) UpdateAll(coll string, selector, update interface{}, opts ...options.UpdateOptions) error {
	return m.db.UpdateAll(coll, selector, update, opts...)
}

func (m *Model) Upsert(coll string, selector, update interface{}, opts ...options.UpsertOptions) error {
	return m.db.Upsert(coll, selector, update, opts...)
}

func (m *Model) Find(coll string, query interface{}, opts ...options.FindOptions) qmgo.QueryI {
	return m.db.Find(coll, query, opts...)
}

func (m *Model) Transaction(f ...func() error) (interface{}, error) {
	return m.db.Transaction(f...)
}
