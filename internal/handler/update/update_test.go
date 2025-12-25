package update

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandlerReturnsBadRequest(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rr := httptest.NewRecorder()

			DefaultHandler(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("DefaultHandler returned status %d, want %d", rr.Code, http.StatusBadRequest)
			}
		})
	}
}
