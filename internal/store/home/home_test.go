package home_test

import (
	"context"
	"testing"

	"github.com/1995parham-teaching/fandogh/internal/config"
	"github.com/1995parham-teaching/fandogh/internal/db"
	"github.com/1995parham-teaching/fandogh/internal/fs"
	"github.com/1995parham-teaching/fandogh/internal/model"
	"github.com/1995parham-teaching/fandogh/internal/store/home"
	"github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

type CommonHomeSuite struct {
	suite.Suite

	Store home.Home
}

func (suite *CommonHomeSuite) TestNoID() {
	require := suite.Require()

	_, err := suite.Store.Get(context.Background(), "invalid_id")
	require.Equal(home.ErrIDNotFound, err)
}

// nolint: funlen
func (suite *CommonHomeSuite) TestSetGet() {
	require := suite.Require()

	cases := []struct {
		name           string
		home           model.Home
		photos         []model.Photo
		expectedSetErr error
		expectedGetErr error
	}{
		{
			name: "Without Error",
			home: model.Home{
				ID:              "",
				Title:           "127.0.0.1",
				Owner:           "parham.alvani@gmail.com",
				Location:        "Iran, Tehran",
				Description:     "Home Sweet Home",
				Peoples:         4,
				Room:            "room_type",
				Bed:             model.Double,
				Rooms:           2,
				Bathrooms:       2,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   true,
				Contract:        "contract_type",
				SecurityDeposit: 0,
				Photos:          nil,
				Price:           0,
			},
			photos: []model.Photo{
				{
					Name:        "1.png",
					ContentType: "image/png",
					Content:     []byte{'1', '2', '3'},
				},
			},
			expectedSetErr: nil,
			expectedGetErr: nil,
		},
		{
			name: "ID Not Empty",
			home: model.Home{
				ID:              "1378",
				Title:           "127.0.0.1",
				Owner:           "parham.alvani@gmail.com",
				Location:        "Iran, Tehran",
				Description:     "Home Sweet Home",
				Peoples:         4,
				Room:            "room_type",
				Bed:             model.Double,
				Rooms:           2,
				Bathrooms:       2,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   true,
				Contract:        "contract_type",
				SecurityDeposit: 0,
				Photos:          nil,
				Price:           0,
			},
			photos: []model.Photo{
				{
					Name:        "1.png",
					ContentType: "image/png",
					Content:     []byte{'1', '2', '3'},
				},
			},
			expectedSetErr: home.ErrIDNotEmpty,
			expectedGetErr: nil,
		},
	}

	for _, c := range cases {
		suite.Run(c.name, func() {
			require.Equal(c.expectedSetErr, suite.Store.Set(context.Background(), &c.home, c.photos))

			if c.expectedSetErr == nil {
				require.NotEmpty(c.home.ID)

				require.NotNil(c.home.Photos)

				for _, key := range c.home.Photos {
					require.NotEmpty(key)
				}

				home, err := suite.Store.Get(context.Background(), c.home.ID)
				require.Equal(c.expectedGetErr, err)

				if c.expectedGetErr == nil {
					require.Equal(c.home, home)
				}
			}
		})
	}
}

type MongoHomeSuite struct {
	CommonHomeSuite

	DB  *mongo.Database
	app *fxtest.App
}

func (suite *MongoHomeSuite) SetupSuite() {
	var (
		database  *mongo.Database
		homeStore home.Home
	)

	suite.app = fxtest.New(
		suite.T(),
		fx.Provide(config.Provide),
		fx.Provide(func(cfg config.Config) db.Config {
			return cfg.Database
		}),
		fx.Provide(func(cfg config.Config) fs.Config {
			return cfg.FileStorage
		}),
		fx.Provide(zap.NewNop),
		fx.Provide(noop.NewTracerProvider().Tracer("")),
		fx.Provide(db.Provide),
		fx.Provide(fs.Provide),
		fx.Provide(
			fx.Annotate(home.Provide, fx.As(new(home.Home))),
		),
		fx.Populate(&database, &homeStore),
	)
	suite.app.RequireStart()

	suite.DB = database
	suite.Store = homeStore
}

func (suite *MongoHomeSuite) TearDownSuite() {
	_, err := suite.DB.Collection(user.Collection).DeleteMany(context.Background(), bson.D{})
	suite.Require().NoError(err)

	suite.app.RequireStop()
}

func TestMongoHomeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MongoHomeSuite))
}
