package qmg

import (
	"context"
	"encoding/json"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Database interface {
		Insert(coll string, docs interface{}, opts ...options.InsertOneOptions) (string, error)
		InsertMany(coll string, docs interface{}, opts ...options.InsertManyOptions) ([]string, error)
		Remove(coll string, selector interface{}, opts ...options.RemoveOptions) error
		RemoveAll(coll string, selector interface{}, opts ...options.RemoveOptions) error
		Update(coll string, selector, update interface{}, opts ...options.UpdateOptions) error
		UpdateAll(coll string, selector, update interface{}, opts ...options.UpdateOptions) error
		Upsert(coll string, selector, update interface{}, opts ...options.UpsertOptions) error
		Find(coll string, query interface{}, opts ...options.FindOptions) qmgo.QueryI
		Transaction(f ...func() error) (interface{}, error)
	}
	decoratedCli struct {
		cli  *qmgo.Client
		ctx  context.Context
		db   *qmgo.Database
		name string
	}
)

func newCli(ctx context.Context, cli *qmgo.Client, db *qmgo.Database) Database {
	return &decoratedCli{
		cli:  cli,
		ctx:  ctx,
		db:   db,
		name: db.GetDatabaseName(),
	}
}

func (c *decoratedCli) Insert(coll string, doc interface{}, opts ...options.InsertOneOptions) (string, error) {
	result, err := c.db.Collection(coll).InsertOne(c.ctx, doc, opts...)
	c.logDuration("insertOne", err, doc)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (c *decoratedCli) InsertMany(coll string, docs interface{}, opts ...options.InsertManyOptions) ([]string, error) {
	result, err := c.db.Collection(coll).InsertMany(c.ctx, docs, opts...)
	c.logDuration("insertMany", err, docs)
	if err != nil {
		return nil, err
	}

	ids := []string{}
	for _, id := range result.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}

	return ids, nil
}

func (c *decoratedCli) Remove(coll string, selector interface{}, opts ...options.RemoveOptions) (err error) {
	err = c.db.Collection(coll).Remove(c.ctx, selector, opts...)
	c.logDuration("remove", err, selector)
	return
}

func (c *decoratedCli) RemoveAll(coll string, selector interface{}, opts ...options.RemoveOptions) (err error) {
	_, err = c.db.Collection(coll).RemoveAll(c.ctx, selector, opts...)
	c.logDuration("removeAll", err, selector)
	return
}

func (c *decoratedCli) Update(coll string, selector, update interface{}, opts ...options.UpdateOptions) (err error) {
	err = c.db.Collection(coll).UpdateOne(c.ctx, selector, update, opts...)
	c.logDuration("updateOne", err, update)
	return
}

func (c *decoratedCli) UpdateAll(coll string, selector, update interface{}, opts ...options.UpdateOptions) (err error) {
	_, err = c.db.Collection(coll).UpdateAll(c.ctx, selector, update, opts...)
	c.logDuration("updateAll", err, update)
	return
}

func (c *decoratedCli) Upsert(coll string, selector, update interface{}, opts ...options.UpsertOptions) (err error) {
	_, err = c.db.Collection(coll).Upsert(c.ctx, selector, update, opts...)
	c.logDuration("upsert", err, update)
	return
}

func (c *decoratedCli) Find(coll string, query interface{}, opts ...options.FindOptions) qmgo.QueryI {
	c.logDuration("find", nil, query)
	return c.db.Collection(coll).Find(c.ctx, query, opts...)
}

func (c *decoratedCli) Transaction(f ...func() error) (interface{}, error) {
	callback := func(sessCtx context.Context) (interface{}, error) {
		for _, v := range f {
			if err := v(); err != nil {
				c.logDuration("transaction", err, nil)
				return nil, err
			}
		}
		return nil, nil
	}
	return c.cli.DoTransaction(c.ctx, callback)
}

func (c *decoratedCli) logDuration(method string, err error, docs ...interface{}) {
	startTime := timex.Now()
	duration := timex.Since(startTime)
	content, e := json.Marshal(docs)
	if e != nil {
		logx.Error(err)
	} else if err != nil {
		if duration > slowThreshold.Load() {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - fail(%s) - %s",
				c.name, method, err.Error(), string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - fail(%s) - %s",
				c.name, method, err.Error(), string(content))
		}
	} else {
		if duration > slowThreshold.Load() {
			logx.WithDuration(duration).Slowf("[MONGO] mongo(%s) - slowcall - %s - ok - %s",
				c.name, method, string(content))
		} else {
			logx.WithDuration(duration).Infof("mongo(%s) - %s - ok - %s", c.name, method, string(content))
		}
	}
}
