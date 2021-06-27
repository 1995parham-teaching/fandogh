package home

import (
	"context"
	"errors"
	"fmt"

	"github.com/1995parham/fandogh/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

var ErrIDNotFound = errors.New("home id does not exist")

// MongoURL communicate with homes collection in MongoDB.
type MongoHome struct {
	DB     *mongo.Database
	Tracer trace.Tracer
}

// Collection is a name of the MongoDB collection for homes.
const Collection = "homes"

// NewMongoHome creates new Home store.
func NewMongoHome(db *mongo.Database, tracer trace.Tracer) *MongoHome {
	return &MongoHome{
		DB:     db,
		Tracer: tracer,
	}
}

// Set saves given home in database and returns its id.
func (s *MongoHome) Set(ctx context.Context, home model.Home) (string, error) {
	ctx, span := s.Tracer.Start(ctx, "store.home.set")
	defer span.End()

	users := s.DB.Collection(Collection)

	result, err := users.InsertOne(ctx, home)
	if err != nil {
		span.RecordError(err)

		return "", fmt.Errorf("mongodb failed: %w", err)
	}

	id, ok := result.InsertedID.(string)
	if !ok {
		panic("invalid mongodb id type")
	}

	return id, nil
}

// Get retrieves home of the given id if it exists.
func (s *MongoHome) Get(ctx context.Context, id string) (model.Home, error) {
	ctx, span := s.Tracer.Start(ctx, "store.home.get")
	defer span.End()

	record := s.DB.Collection(Collection).FindOne(ctx, bson.M{
		"_id": id,
	})

	var home model.Home
	if err := record.Decode(&home); err != nil {
		span.RecordError(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return home, ErrIDNotFound
		}

		return home, fmt.Errorf("mongodb failed: %w", err)
	}

	return home, nil
}
