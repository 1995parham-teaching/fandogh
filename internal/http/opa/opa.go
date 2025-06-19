package opa

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/open-policy-agent/opa/v1/sdk"
	sdktest "github.com/open-policy-agent/opa/v1/sdk/test"

	"github.com/1995parham-teaching/fandogh/internal/http/common"
)

type OPA struct {
	engine *sdk.OPA
}

func New() (OPA, error) {
	policyBytes, err := os.ReadFile("policy.rego")
	if err != nil {
		return OPA{ engine: nil}, fmt.Errorf("failed to read policy file: %w", err)
	}

	policy := string(policyBytes)

	// create a mock HTTP bundle server (in production it should be on its server)
	server, err := sdktest.NewServer(sdktest.MockBundle("/bundles/bundle.tar.gz", map[string]string{
		"example.rego": policy,
	}))
	if err != nil {
		return OPA{ engine: nil}, fmt.Errorf("failed to run opa test server: %w", err)
	}

	// provide the OPA configuration which specifies
	// fetching policy bundles from the mock server
	// and logging decisions locally to the console
	config := fmt.Sprintf(`{
		"services": {
			"test": {
				"url": %q
			}
		},
		"bundles": {
			"test": {
				"resource": "/bundles/bundle.tar.gz"
			}
		},
		"decision_logs": {
			"console": true
		}
	}`, server.URL())

	eng, err := sdk.New(context.Background(), sdk.Options{ // nolint: exhaustruct
		ID:     "opa-test-1",
		Config: strings.NewReader(config),
		Ready: make(chan struct{}),
	})
	if err != nil {
		return OPA{
			engine: nil,
		}, fmt.Errorf("failed to create opa engine %w", err)
	}

	return OPA{
		engine: eng,
	}, nil
}

func (opa OPA) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mc, ok := c.Get(common.UserContextKey).(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, "user claims not found")
			}

			sub, err := mc.Claims.GetSubject()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "user claims subject not found")
			}

			result, err := opa.engine.Decision(
				c.Request().Context(),
				sdk.DecisionOptions{ // nolint: exhaustruct
					Path: "/authz/allow",
					Input: map[string]any{
						"method": c.Request().Method,
						"subject": sub,
						"path": c.Path(),
					},
				},
			)
			if err != nil {
				return echo.ErrForbidden
			}

			if decision, ok := result.Result.(bool); !ok || !decision {
				return echo.NewHTTPError(http.StatusForbidden, "policy validation failed")
			}

			return next(c)
		}
	}
}
