package update

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandlerReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	DefaultHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("DefaultHandler returned status %d, want %d", status, http.StatusBadRequest)
	}
}

