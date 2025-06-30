package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/1995parham-teaching/fandogh/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrEmailNotFound indicates that given email does not exist on database.
	ErrEmailNotFound = errors.New("given email does not exist")
	// ErrEmailDuplicate indicates that given email is exists on database.
	ErrEmailDuplicate = errors.New("given email exists")
)

// MongoUser communicate with users collection in MongoDB.
type MongoUser struct {
	DB     *mongo.Database
	Tracer trace.Tracer
}

// Collection is a name of the MongoDB collection for Users.
const Collection = "users"

// NewMongoUser creates new User store.
func NewMongoUser(db *mongo.Database, tracer trace.Tracer) *MongoUser {
	return &MongoUser{
		DB:     db,
		Tracer: tracer,
	}
}

// Set saves given user in database.
func (s *MongoUser) Set(ctx context.Context, user model.User) error {
	ctx, span := s.Tracer.Start(ctx, "store.user.set")
	defer span.End()

	users := s.DB.Collection(Collection)

	_, err := users.InsertOne(ctx, user)
	if err != nil {
		span.RecordError(err)

		if mongo.IsDuplicateKeyError(err) {
			return ErrEmailDuplicate
		}

		return fmt.Errorf("mongodb failed: %w", err)
	}

	return nil
}

// Get retrieves user of the given email if it exists.
func (s *MongoUser) Get(ctx context.Context, email string) (model.User, error) {
	ctx, span := s.Tracer.Start(ctx, "store.url.get")
	defer span.End()

	record := s.DB.Collection(Collection).FindOne(ctx, bson.M{
		"email": email,
	})

	var user model.User

	err := record.Decode(&user)
	if err != nil {
		span.RecordError(err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, ErrEmailNotFound
		}

		return user, fmt.Errorf("mongodb failed: %w", err)
	}

	return user, nil
}
