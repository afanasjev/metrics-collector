package update

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	rr := httptest.NewRecorder()

	DefaultHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("DefaultHandler() status = %v, want %v", status, http.StatusBadRequest)
	}
}

