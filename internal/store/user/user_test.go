package user_test

import (
	"context"
	"testing"

	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/model"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/trace/noop"
)

type CommonUserSuite struct {
	suite.Suite

	Store user.User
}

func (suite *CommonUserSuite) TestNoEmail() {
	require := suite.Require()

	_, err := suite.Store.Get(context.Background(), "notexists@gmail.com")
	require.Equal(user.ErrEmailNotFound, err)
}

func (suite *CommonUserSuite) TestSetGet() {
	require := suite.Require()

	cases := []struct {
		name           string
		user           model.User
		expectedSetErr error
		expectedGetErr error
	}{
		{
			name: "Without Error",
			user: model.User{
				Name:     "Parham Alvani",
				Email:    "parham.alvani@gmail.com",
				Password: "123456",
			},
			expectedSetErr: nil,
			expectedGetErr: nil,
		},
		{
			name: "Duplicate Key",
			user: model.User{
				Name:     "Parham Alvani",
				Email:    "parham.alvani@gmail.com",
				Password: "123456",
			},
			expectedSetErr: user.ErrEmailDuplicate,
			expectedGetErr: nil,
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			require.Equal(c.expectedSetErr, suite.Store.Set(context.Background(), c.user))

			if c.expectedSetErr == nil {
				user, err := suite.Store.Get(context.Background(), c.user.Email)
				require.Equal(c.expectedGetErr, err)

				if c.expectedGetErr == nil {
					require.Equal(c.user, user)
				}
			}
		})
	}
}

type MongoUserSuite struct {
	CommonUserSuite

	DB *mongo.Database
}

func (suite *MongoUserSuite) SetupSuite() {
	cfg := config.New()

	db, err := db.New(cfg.Database)
	suite.Require().NoError(err)

	suite.DB = db
	suite.Store = user.NewMongoUser(db, noop.NewTracerProvider().Tracer(""))
}

func (suite *MongoUserSuite) TearDownSuite() {
	_, err := suite.DB.Collection(user.Collection).DeleteMany(context.Background(), bson.D{})
	suite.Require().NoError(err)

	suite.Require().NoError(suite.DB.Client().Disconnect(context.Background()))
}

func TestMongoUserSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MongoUserSuite))
}

type MemoryUserSuite struct {
	CommonUserSuite
}

func (suite *MemoryUserSuite) SetupSuite() {
	suite.Store = user.NewMemoryUser()
}

func TestMemoryUserSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MemoryUserSuite))
}
