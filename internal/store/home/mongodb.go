package home

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/1995parham-teaching/fandogh/internal/fs"
	"github.com/1995parham-teaching/fandogh/internal/model"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrIDNotFound = errors.New("home id does not exist")
	ErrIDNotEmpty = errors.New("home id must be empty")
)

// MongoHome communicate with homes collection in MongoDB.
type MongoHome struct {
	DB     *mongo.Database
	Minio  *minio.Client
	Tracer trace.Tracer
}

const (
	// Collection is a name of the MongoDB collection for homes.
	Collection = "homes"

	// Bucket for storing photos.
	Bucket = "photos"
)

// NewMongoHome creates new Home store.
func NewMongoHome(db *mongo.Database, client *minio.Client, tracer trace.Tracer) *MongoHome {
	return &MongoHome{
		DB:     db,
		Tracer: tracer,
		Minio:  client,
	}
}

// Set saves given home in database and returns its id.
func (s *MongoHome) Set(ctx context.Context, home *model.Home, photos []model.Photo) error {
	ctx, span := s.Tracer.Start(ctx, "store.home.set")
	defer span.End()

	err := fs.Bucket(ctx, s.Minio, Bucket)
	if err != nil {
		span.RecordError(ErrIDNotEmpty)

		return fmt.Errorf("minio bucket creation/checking failed: %w", err)
	}

	if home.ID != "" {
		span.RecordError(ErrIDNotEmpty)

		return ErrIDNotEmpty
	}

	home.ID = bson.NewObjectID().Hex()

	if home.Photos == nil {
		home.Photos = make(map[string]string)
	}

	for _, photo := range photos {
		home.Photos[photo.Name] = fs.Generate(home.ID, photo.Name)

		// nolint: exhaustruct
		_, err = s.Minio.PutObject(ctx, Bucket, home.Photos[photo.Name],
			bytes.NewReader(photo.Content), int64(len(photo.Content)), minio.PutObjectOptions{
				ContentType: photo.ContentType,
			})
		if err != nil {
			span.RecordError(err)

			return fmt.Errorf("minio object creation failed: %w", err)
		}
	}

	users := s.DB.Collection(Collection)

	_, err = users.InsertOne(ctx, home)
	if err != nil {
		span.RecordError(err)

		return fmt.Errorf("mongodb failed: %w", err)
	}

	return nil
}

// Get retrieves home of the given id if it exists.
func (s *MongoHome) Get(ctx context.Context, id string) (model.Home, error) {
	ctx, span := s.Tracer.Start(ctx, "store.home.get")
	defer span.End()

	record := s.DB.Collection(Collection).FindOne(ctx, bson.M{
		"_id": id,
	})

	var home model.Home

	err := record.Decode(&home)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return home, ErrIDNotFound
		}

		return home, fmt.Errorf("mongodb failed: %w", err)
	}

	return home, nil
}
