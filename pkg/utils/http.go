package utils

import (
	"context"
	"encoding/json"
	"io"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/kkr2/spots/pkg/httpErrors"
	"github.com/kkr2/spots/pkg/logger"
	"github.com/kkr2/spots/pkg/sanitize"
)

// GetRequestID gets id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// ReqIDCtxKey is a key used for the Request ID in context
type ReqIDCtxKey struct{}

// GetCtxWithReqID timeout and request id from echo context
func GetCtxWithReqID(c echo.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*15)
	ctx = context.WithValue(ctx, ReqIDCtxKey{}, GetRequestID(c))
	return ctx, cancel
}

// GetRequestCtx gets context from request
func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

// GetConfigPath gets path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config-docker"
	}
	return "./internal/config/config-local"
}

// GetIPAddress gets user ip address
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// ErrResponseWithLog response with logging error for echo context
func ErrResponseWithLog(ctx echo.Context, logger logger.Logger, err error) error {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
	return ctx.JSON(httpErrors.ErrorResponse(err))
}

// LogResponseError response with logging error for echo context
func LogResponseError(ctx echo.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
}

// ReadRequest request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

// SanitizeRequest sanitize and validate request
func SanitizeRequest(ctx echo.Context, request interface{}) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	defer ctx.Request().Body.Close()

	sanBody, err := sanitize.SanitizeJSON(body)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err = json.Unmarshal(sanBody, request); err != nil {
		return err
	}

	return validate.StructCtx(ctx.Request().Context(), request)
}
