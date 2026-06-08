package mongo

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type QueryOptions struct {
	Filter bson.M
	Sort   bson.D
	Limit  int
	Offset int
}

type QueryOption func(*QueryOptions)

func WithFilter(filter bson.M) QueryOption {
	return func(q *QueryOptions) {
		q.Filter = filter
	}
}

func WithSort(field string, asc bool) QueryOption {
	dir := -1
	if asc {
		dir = 1
	}
	return func(q *QueryOptions) {
		q.Sort = append(q.Sort, bson.E{Key: field, Value: dir})
	}
}

func WithLimit(limit int) QueryOption {
	return func(q *QueryOptions) {
		q.Limit = limit
	}
}

func WithOffset(offset int) QueryOption {
	return func(q *QueryOptions) {
		q.Offset = offset
	}
}

func buildQueryOptions(opts []QueryOption) *QueryOptions {
	q := &QueryOptions{
		Filter: bson.M{},
	}
	for _, o := range opts {
		o(q)
	}
	return q
}

func (q *QueryOptions) toFindOptions() *options.FindOptionsBuilder {
	opt := options.Find()
	if len(q.Sort) > 0 {
		opt.SetSort(q.Sort)
	}
	if q.Limit > 0 {
		opt.SetLimit(int64(q.Limit))
	}
	if q.Offset > 0 {
		opt.SetSkip(int64(q.Offset))
	}
	return opt
}
