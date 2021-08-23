package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServerCtxKeys string

// Context values from GraphQL Server
const (
	GIN_CTX_KEY    ServerCtxKeys = "gin_context_key"
	REQUEST_ID_KEY ServerCtxKeys = "request_id_key"
	OP_TYPE_KEY    ServerCtxKeys = "operation_type_key"
	OP_NAME_KEY    ServerCtxKeys = "operation_name_key"
	OP_RAW         ServerCtxKeys = "operation_raw_key"
)

// LogValues represents the values we want to include in server logs
type LogValues struct {
	ReqId      string
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	ForwardFor []string
	MsgStr     string
	Type       interface{}
	Action     interface{}
	Body       interface{}
}

// logValuesFromCtx values from the context ready to be logged
func logValuesFromCtx(serName string, c *gin.Context) *LogValues {
	t := time.Now()

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	if raw != "" {
		path = path + "?" + raw
	}
	msg := c.Errors.String()
	if msg == "" {
		msg = "Request"
	}

	reqId := c.Request.Context().Value(REQUEST_ID_KEY).(string)
	actionType := c.Request.Context().Value(OP_TYPE_KEY)
	action := c.Request.Context().Value(OP_NAME_KEY)
	body := c.Request.Context().Value(OP_RAW)
	return &LogValues{
		ReqId:      reqId,
		SerName:    serName,
		Latency:    time.Since(t),
		StatusCode: c.Writer.Status(),
		ClientIP:   c.ClientIP(),
		ForwardFor: c.Request.Header.Values("X-Forwarded-For"),
		MsgStr:     msg,
		Type:       actionType,
		Action:     action,
		Body:       body,
	}
}

// logSwitch prints LogValues into the right log stream using zerolog
func logSwitch(data *LogValues) {
	var logger *zerolog.Event
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		logger = log.Warn()
	case data.StatusCode >= 500:
		logger = log.Error()
	default:
		logger = log.Info()
	}

	logger.
		Str("req_id", data.ReqId).
		Interface("req_body", data.Body).
		Str("ser_name", data.SerName).
		Dur("resp_time", data.Latency).
		Int("status", data.StatusCode).
		Str("client_ip", data.ClientIP).
		Strs("fwd_for", data.ForwardFor).
		Interface("op_name", data.Action).
		Interface("op_type", data.Type).
		Msg(data.MsgStr)
}

// LogMiddleware is a Gin server logger middleware
func LogMiddleware(serName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logSwitch(logValuesFromCtx(serName, c))
	}
}

// ErrorLogger call this function to log errors from a GraphQL context
func ErrorLogger(ctx context.Context) {
	// TODO: Implement bug tracker
	gc, err := GinCtxFromCtx(ctx)
	if err == nil {
		logSwitch(logValuesFromCtx("gin", gc))
	} else {
		log.Error().Msgf("Unexpected error: %v", err)
	}
}

// GinCtxToCtxMiddleware stores Gin request context into a generic context
func GinCtxToCtxMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GIN_CTX_KEY, c)
		ctx = context.WithValue(ctx, REQUEST_ID_KEY, "Franklin was here")

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinCtxFromCtx extracts Gin request context from a generic context
// if the Gin context is not present returns an error
func GinCtxFromCtx(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GIN_CTX_KEY)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

// GQLOperationMiddleware can be passed into Gqlgen server as AroundOperation middleware
func GQLOperationMiddleware(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	oc := graphql.GetOperationContext(ctx)
	gc, err := GinCtxFromCtx(ctx)
	if err == nil {
		req := gc.Request.Context()
		req = context.WithValue(req, OP_TYPE_KEY, oc.Operation.Operation)
		req = context.WithValue(req, OP_NAME_KEY, oc.OperationName)
		req = context.WithValue(req, OP_RAW, oc.RawQuery)
		// TODO: Should we log operation stats?
		gc.Request = gc.Request.WithContext(req)

		newCtx := context.WithValue(ctx, GIN_CTX_KEY, gc)
		return next(newCtx)
	}

	return next(ctx)
}
