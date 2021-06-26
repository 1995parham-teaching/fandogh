package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1995parham/fandogh/internal/http/handler"
	"github.com/1995parham/fandogh/internal/http/request"
	store "github.com/1995parham/fandogh/internal/store/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserSuite struct {
	suite.Suite

	engine *echo.Echo
}

func (suite *UserSuite) SetupSuite() {
	suite.engine = echo.New()

	user := handler.User{
		Store:  store.NewMemoryUser(),
		Logger: zap.NewNop(),
		Tracer: trace.NewNoopTracerProvider().Tracer(""),
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

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(b))

	suite.engine.ServeHTTP(w, req)
	require.Equal(http.StatusBadRequest, w.Code)
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
			req := httptest.NewRequest("POST", "/register", bytes.NewReader(b))
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
