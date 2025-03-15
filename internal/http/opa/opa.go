package opa

import (
	"net/http"

	"github.com/1995parham-teaching/fandogh/internal/http/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/open-policy-agent/opa/v1/sdk"
)

type OPA struct {
	Engine *sdk.OPA
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

			allow, err := opa.Engine.Decision(
				c.Request().Context(),
				sdk.DecisionOptions{
					Path: "/authz/allow",
					Input: map[string]any{
						"subject": sub,
					},
				},
			)
			if err != nil {
				return echo.ErrForbidden
			}

			obj, ok := allow.Result.(map[string]any)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "policy validation failed")
			}

			if access, ok := obj["allow"].(bool); !ok || !access {
				if body, ok := obj["body"].(string); ok {
					return echo.NewHTTPError(http.StatusForbidden, body)
				}

				return echo.NewHTTPError(http.StatusForbidden, "policy validation failed")
			}

			return next(c)
		}
	}
}
