package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anishsharma21/go-web-dev-template/internal"
	"github.com/google/uuid"
)

// LoggingMiddleware logs requests and responses
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), internal.REQUEST_ID_KEY, requestID)
		r = r.WithContext(ctx)

		// Wrap response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Milliseconds()

		logLevel := slog.LevelInfo
		if rw.statusCode >= 400 && rw.statusCode < 500 {
			logLevel = slog.LevelWarn
		} else if rw.statusCode >= 500 {
			logLevel = slog.LevelError
		}

		// Skip logging for favicon requests
		if !strings.Contains(r.URL.Path, "favicon") {
			slog.Log(ctx, logLevel, fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, rw.statusCode),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status_code", rw.statusCode),
				slog.Int64("processing_ms", duration),
				slog.String("client_ip", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
		}
	})
}

// CustomLogHandler adds request_id and user_id fields to log records
type CustomLogHandler struct {
	slog.Handler
}

// Handle adds request_id and user_id fields to the log record if they exist in the context
func (h *CustomLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(internal.REQUEST_ID_KEY).(string); ok {
		r.AddAttrs(slog.String(internal.REQUEST_ID_KEY, requestID))
	}
	if userID, ok := ctx.Value(internal.CLERK_USER_ID_KEY).(string); ok {
		r.AddAttrs(slog.String(internal.CLERK_USER_ID_KEY, userID))
	}
	return h.Handler.Handle(ctx, r)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
