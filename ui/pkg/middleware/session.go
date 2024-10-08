package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dskart/waterfall-engine/app"
	uiErrors "github.com/dskart/waterfall-engine/ui/pkg/errors"
	"go.uber.org/zap"
)

type contextKey string

const sessionContextKey contextKey = "SessionContextKey"

func CtxSession(ctx context.Context) *app.Session {
	return ctx.Value(sessionContextKey).(*app.Session).WithContext(ctx)
}

const MaxRequestBytes int64 = 20 * 1024 * 1024

type SessionMiddleware interface {
	AnonymousSession(next http.Handler) http.Handler
}

// Define our struct
type sessionMiddleware struct {
	app        *app.App
	proxyCount int
}

func (sm *sessionMiddleware) AnonymousSession(f http.Handler) http.Handler {
	return sm.createSession(f)
}

func (sm *sessionMiddleware) createSession(f http.Handler) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		beginTime := time.Now()
		session := sm.app.NewSession(sm.app.Logger()).WithContext(r.Context())

		hijacker, _ := w.(http.Hijacker)
		w = &uiErrors.StatusCodeRecorder{
			ResponseWriter: w,
			Hijacker:       hijacker,
		}

		id := make([]byte, 20)
		if _, err := rand.Read(id); err != nil {
			session.Logger().Error("failed to generate random request ID: %s", zap.Error(err))
			http.Error(w, "failed to generate ID", http.StatusInternalServerError)
			return
		}
		requestId := base64.RawURLEncoding.EncodeToString(id)
		remote := httpRequestIPAddress(r, sm.proxyCount)
		session = session.WithLogger(session.Logger().With(
			zap.String("requestId", requestId),
			zap.String("peer.address", remote),
		))

		defer func() {
			statusCode := w.(*uiErrors.StatusCodeRecorder).StatusCode
			if statusCode == 0 {
				statusCode = 200
			}
			duration := time.Since(beginTime)
			logger := session.Logger().With(
				zap.Duration("duration", duration),
				zap.Uint32("status_code", uint32(statusCode)),
			)
			url := *r.URL
			// Scrub query values
			url.RawQuery = ""
			logger.Info(r.Method + " " + url.RequestURI())
		}()

		defer func() {
			if r := recover(); r != nil {
				session.Logger().Sugar().Errorf("%v: %s", r, debug.Stack())
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBytes)

		r = r.WithContext(context.WithValue(r.Context(), sessionContextKey, session))
		f.ServeHTTP(w, r)
	}

	return http.HandlerFunc(handlerFunc)
}

func NewSessionMiddleware(app *app.App, proxyCount int) SessionMiddleware {
	return &sessionMiddleware{
		app:        app,
		proxyCount: proxyCount,
	}
}

func httpRequestIPAddress(r *http.Request, proxyCount int) string {
	addr := r.RemoteAddr
	if h := r.Header.Get("X-Forwarded-For"); h != "" && proxyCount > 0 {
		if clients := strings.Split(h, ","); proxyCount > len(clients) {
			addr = clients[0]
		} else {
			addr = clients[len(clients)-proxyCount]
		}
	}
	addr = strings.TrimSpace(addr)
	if host, _, err := net.SplitHostPort(addr); err == nil && host != "" {
		return host
	}
	return addr
}
