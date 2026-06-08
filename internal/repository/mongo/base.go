package mongo

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Document interface {
	GetID() bson.ObjectID
	SetID(id bson.ObjectID)
	SetTimestamnps(createdAt time.Time, updatedAt time.Time)
}

type BaseRepository[T Document] struct {
	col *mongo.Collection
}

func NewBaseRepository[T Document](db *mongo.Database, collectionName string) *BaseRepository[T] {
	return &BaseRepository[T]{col: db.Collection(collectionName)}
}

func (r *BaseRepository[T]) Insert(ctx context.Context, doc T) error {
	doc.SetID(bson.NewObjectID())
	now := time.Now()
	doc.SetTimestamnps(now, now)

	_, err := r.col.InsertOne(ctx, doc)
	return err
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, opts ...QueryOption) (T, error) {
	q := buildQueryOptions(opts)

	var result T
	err := r.col.FindOne(ctx, q.Filter).Decode(&result)

	return result, err
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id string) (T, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		var zero T
		return zero, err
	}
	return r.FindOne(ctx, WithFilter(bson.M{"_id": oid}))
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, opts ...QueryOption) ([]T, error) {
	q := buildQueryOptions(opts)

	cursor, err := r.col.Find(ctx, q.Filter, q.toFindOptions())
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			slog.Error("failed to close cursor", "error", cerr)
		}
	}()

	var result []T
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *BaseRepository[T]) Count(ctx context.Context, opts ...QueryOption) (int64, error) {
	q := buildQueryOptions(opts)
	return r.col.CountDocuments(ctx, q.Filter)
}

func (r *BaseRepository[T]) UpdateOne(ctx context.Context, id string, update bson.M) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if set, ok := update["$set"].(bson.M); ok {
		set["updated_at"] = time.Now()
	}
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

func (r *BaseRepository[T]) DeleteOne(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
