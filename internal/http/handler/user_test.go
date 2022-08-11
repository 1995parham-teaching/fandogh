package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1995parham-teaching/fandogh/internal/http/handler"
	"github.com/1995parham-teaching/fandogh/internal/http/jwt"
	"github.com/1995parham-teaching/fandogh/internal/http/request"
	"github.com/1995parham-teaching/fandogh/internal/model"
	store "github.com/1995parham-teaching/fandogh/internal/store/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserSuite struct {
	suite.Suite

	store  store.User
	engine *echo.Echo
}

func (suite *UserSuite) SetupSuite() {
	suite.engine = echo.New()
	suite.store = store.NewMemoryUser()

	user := handler.User{
		Store:  suite.store,
		Logger: zap.NewNop(),
		Tracer: trace.NewNoopTracerProvider().Tracer(""),
		JWT: jwt.JWT{
			Config: jwt.Config{
				AccessTokenSecret: "secret",
			},
		},
	}

	user.Register(suite.engine.Group(""))
}

func (suite *UserSuite) TestBadRequest() {
	require := suite.Require()

	// because there is no content-type header, request is categorized as a bad request.
	b, err := json.Marshal(request.Register{
		Name:     "Parham Alvani",
		Password: "123456",
		Email:    "parham.alvani@gmail.com",
	})
	require.NoError(err)

	{
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))

		suite.engine.ServeHTTP(w, req)
		require.Equal(http.StatusBadRequest, w.Code)
	}
	{
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))

		suite.engine.ServeHTTP(w, req)
		require.Equal(http.StatusBadRequest, w.Code)
	}
}

func (suite *UserSuite) TestRegister() {
	require := suite.Require()

	cases := []struct {
		name     string
		code     int
		register request.Register
	}{
		{
			name: "Successful",
			code: http.StatusCreated,
			register: request.Register{
				Name:     "Parham Alvani",
				Email:    "parham.alvani@gmail.com",
				Password: "123456",
			},
		}, {
			name: "Duplicate Key",
			code: http.StatusBadRequest,
			register: request.Register{
				Name:     "Parham Alvani",
				Email:    "parham.alvani@gmail.com",
				Password: "123456",
			},
		}, {
			name: "Invalid URL",
			code: http.StatusBadRequest,
			register: request.Register{
				Name:     "Parham Alvani",
				Email:    "parham.alvani@gmail",
				Password: "123456",
			},
		},
	}

	for _, c := range cases {
		c := c
		suite.Run(c.name, func() {
			b, err := json.Marshal(c.register)
			require.NoError(err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			require.Equal(c.code, w.Code)
		})
	}
}

func (suite *UserSuite) TestLogin() {
	require := suite.Require()

	require.NoError(suite.store.Set(context.Background(), model.User{
		Name:     "Elahe Dastan",
		Email:    "elahe.dstn@gmail.com",
		Password: "123456",
	}))

	cases := []struct {
		name  string
		code  int
		login request.Login
	}{
		{
			name: "Successful",
			code: http.StatusOK,
			login: request.Login{
				Email:    "elahe.dstn@gmail.com",
				Password: "123456",
			},
		}, {
			name: "Incorrect Password",
			code: http.StatusUnauthorized,
			login: request.Login{
				Email:    "elahe.dstn@gmail.com",
				Password: "1234567",
			},
		}, {
			name: "No Email",
			code: http.StatusNotFound,
			login: request.Login{
				Email:    "noone@gmail.com",
				Password: "123456",
			},
		}, {
			name: "Inavlid Email",
			code: http.StatusBadRequest,
			login: request.Login{
				Email:    "@gmail.com",
				Password: "123456",
			},
		},
	}

	for _, c := range cases {
		c := c
		suite.Run(c.name, func() {
			b, err := json.Marshal(c.login)
			require.NoError(err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			require.Equal(c.code, w.Code)
		})
	}
}

func TestURLSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UserSuite))
}
