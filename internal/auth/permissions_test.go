package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestRequireRole(t *testing.T) {
	protected := RequireRole("admin", "teacher")(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("allows request for allowed role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resource", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "admin"))
		rr := httptest.NewRecorder()

		protected.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
		}
	})

	t.Run("rejects when role missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resource", nil)
		rr := httptest.NewRecorder()

		protected.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}
	})

	t.Run("rejects when role not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resource", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "student"))
		rr := httptest.NewRecorder()

		protected.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})
}

func TestRequireSelfOrAdmin(t *testing.T) {
	protected := RequireSelfOrAdmin()(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("allows admin for any target", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/target-id", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "admin"))
		req = req.WithContext(context.WithValue(req.Context(), "userID", "someone-else"))
		req = mux.SetURLVars(req, map[string]string{"id": "target-id"})
		rr := httptest.NewRecorder()

		protected.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
		}
	})

	t.Run("allows user for own id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/u-1", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "student"))
		req = req.WithContext(context.WithValue(req.Context(), "userID", "u-1"))
		req = mux.SetURLVars(req, map[string]string{"id": "u-1"})
		rr := httptest.NewRecorder()

		protected.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
		}
	})

	t.Run("rejects when not admin and not self", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/u-2", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "student"))
		req = req.WithContext(context.WithValue(req.Context(), "userID", "u-1"))
		req = mux.SetURLVars(req, map[string]string{"id": "u-2"})
		rr := httptest.NewRecorder()

		protected.ServeHTTP(rr, req)

		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
		}
	})
}
