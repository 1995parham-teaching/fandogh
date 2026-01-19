package user_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"

	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/model"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
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
			u := c.user
			require.Equal(c.expectedSetErr, suite.Store.Set(context.Background(), &u))

			if c.expectedSetErr == nil {
				user, err := suite.Store.Get(context.Background(), c.user.Email)
				require.Equal(c.expectedGetErr, err)

				if c.expectedGetErr == nil {
					require.Equal(u.Email, user.Email)
					require.Equal(u.Name, user.Name)
					require.Equal(u.Password, user.Password)
				}
			}
		})
	}
}

type MongoUserSuite struct {
	CommonUserSuite

	DB  *mongo.Database
	app *fxtest.App
}

func (suite *MongoUserSuite) SetupSuite() {
	var (
		database  *mongo.Database
		userStore user.User
	)

	suite.app = fxtest.New(
		suite.T(),
		fx.Provide(config.Provide),
		fx.Provide(zap.NewNop),
		fx.Provide(func() trace.Tracer {
			return noop.NewTracerProvider().Tracer("")
		}),
		fx.Provide(db.Provide),
		fx.Provide(
			fx.Annotate(user.Provide, fx.As(new(user.User))),
		),
		fx.Populate(&database, &userStore),
	)
	suite.app.RequireStart()

	suite.DB = database
	suite.Store = userStore
}

func (suite *MongoUserSuite) TearDownSuite() {
	_, err := suite.DB.Collection(user.Collection).DeleteMany(context.Background(), bson.D{})
	suite.Require().NoError(err)

	suite.app.RequireStop()
}

func TestMongoUserSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MongoUserSuite))
}

type MemoryUserSuite struct {
	CommonUserSuite
}

func (suite *MemoryUserSuite) SetupSuite() {
	var userStore user.User

	app := fxtest.New(
		suite.T(),
		fx.Provide(
			fx.Annotate(
				user.NewMemoryUser,
				fx.As(new(user.User)),
			),
		),
		fx.Populate(&userStore),
	)
	defer app.RequireStart().RequireStop()

	suite.Store = userStore
}

func TestMemoryUserSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MemoryUserSuite))
}
