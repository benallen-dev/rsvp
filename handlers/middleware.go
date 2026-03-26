package handlers

import (
	"context"
	"fmt"
	"net/http"

	"rsvp/storage"
)

type contextKey string

const storeKey contextKey = "store"

// WithStore adds the store to the request context
func WithStore(ctx context.Context, s *storage.Store) context.Context {
	return context.WithValue(ctx, storeKey, s)
}

// StoreFromContext retrieves the store from the request context
func StoreFromContext(ctx context.Context) (*storage.Store, error) {
	s, ok := ctx.Value(storeKey).(*storage.Store)
	if !ok {
		return nil, fmt.Errorf("store not found in context")
	}
	return s, nil
}

// storeMiddleware injects the store into the request context
func storeMiddleware(s *storage.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithStore(r.Context(), s)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// wrapHandler applies middleware to a handler function
func wrapHandler(middleware func(http.Handler) http.Handler, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// noCacheHandler applies no-cache headers to a handler function
func noCacheHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		handler(w, r)
	}
}
