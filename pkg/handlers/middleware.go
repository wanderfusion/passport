package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	requestIdFormat string = "req-%s"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// since we are using requestID only for tracing and not indexing, we can use a string here
		reqId := fmt.Sprintf(requestIdFormat, uuid.New().String())
		ctx := context.WithValue(r.Context(), RequestIdContextKey, reqId)

		log.Info().
			Str("requestId", reqId).
			Str("method", r.Method).
			Str("url", r.URL.RawPath).
			Str("remoteAddress", r.RemoteAddr).
			Str("userAgent", r.UserAgent()).
			Str("protocolVersion", r.Proto).
			Msg("Recieved new request.")

		metrics := httpsnoop.CaptureMetrics(next, w, r.WithContext(ctx))

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		log.Info().
			Str("requestId", reqId).
			Int64("latencyMicroseconds", metrics.Duration.Microseconds()).
			Int("statusCode", metrics.Code).
			Str("route", routePattern).
			Msg("Finished processing request.")
	})
}
