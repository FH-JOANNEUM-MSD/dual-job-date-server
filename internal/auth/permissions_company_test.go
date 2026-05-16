package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Covers GET /api/companies/{id}/meetings RBAC wrapper (same entity type as other company routes).
func TestRequireSelfOrAdmin_company(t *testing.T) {
	protected := RequireSelfOrAdmin("company")(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("admin passes for any numeric company id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/companies/99/meetings", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "admin"))
		req = mux.SetURLVars(req, map[string]string{"id": "99"})
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected %d got %d", http.StatusNoContent, rr.Code)
		}
	})

	t.Run("company passes when URL id matches context company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/companies/7/meetings", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "company"))
		req = req.WithContext(context.WithValue(req.Context(), "company_id", 7))
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected %d got %d", http.StatusNoContent, rr.Code)
		}
	})

	t.Run("company forbidden when URL id differs from context company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/companies/8/meetings", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "company"))
		req = req.WithContext(context.WithValue(req.Context(), "company_id", 7))
		req = mux.SetURLVars(req, map[string]string{"id": "8"})
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("student forbidden for company entity", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/companies/1/meetings", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "student"))
		req = req.WithContext(context.WithValue(req.Context(), "student_id", 1))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("expected %d got %d", http.StatusForbidden, rr.Code)
		}
	})

	t.Run("company forbidden on invalid numeric id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/companies/not-a-slot/meetings", nil)
		req = req.WithContext(context.WithValue(req.Context(), "role", "company"))
		req = req.WithContext(context.WithValue(req.Context(), "company_id", 7))
		req = mux.SetURLVars(req, map[string]string{"id": "not-a-slot"})
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
