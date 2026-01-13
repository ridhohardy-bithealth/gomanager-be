package middleware

import (
	"net/http"

	customErrors "ps-gogo-manajer/pkg/custom-errors"
	jwt "ps-gogo-manajer/pkg/jwt"
	"ps-gogo-manajer/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			jwtToken, err := extractJWTTokenFromHeader(ctx.Request())
			if err != nil {
				return ctx.JSON(response.WriteErrorResponse(err))
			}

			claim, err := jwt.ClaimToken(jwtToken)
			if err != nil {
				err = errors.Wrap(customErrors.ErrUnauthorized, err.Error())
				return ctx.JSON(response.WriteErrorResponse(err))
			}

			ctx.Set("user", claim)

			// default user passing middleware if token is valid
			return next(ctx)
		}
	}
}

func extractJWTTokenFromHeader(r *http.Request) (string, error) {
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return "", errors.Wrap(customErrors.ErrUnauthorized, "missing auth token")
	}

	return authToken[len("Bearer "):], nil
}
