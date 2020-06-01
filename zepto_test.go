package zepto

import (
	"encoding/json"
	"github.com/go-zepto/zepto/testutils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewZepto(t *testing.T) {
	z := NewZepto()
	if z.logger == nil {
		t.Errorf("Logger should not be nil")
	}
}

func TestSetupHTTP(t *testing.T) {
	z := NewZepto()
	h := mux.NewRouter()
	z.SetupHTTP("0.0.0.0:8000", h)
	w := httptest.NewRecorder()
	if z.httpServer == nil {
		t.Fatal("z.httpServer should not be nil")
	}
	if z.httpServer.Handler == nil {
		t.Fatal("z.httpServer.Handler should not be nil")
	}
	now := time.Now()
	z.startedAt = &now
	z.httpServer.Handler.ServeHTTP(w, httptest.NewRequest("GET", "/health", nil))
	var s HealthStatus
	err := json.Unmarshal(w.Body.Bytes(), &s)
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != "zepto" {
		t.Errorf("Expected name to be equal 'zepto'")
	}
	if s.Version != "latest" {
		t.Errorf("Expected version to be equal 'latest'")
	}
	if s.Age != "0s" {
		t.Errorf("Expected age to be equal '0s'")
	}
}

func TestSetupGRPC(t *testing.T) {
	z := NewZepto()
	z.SetupGRPC("0.0.0.0:9000", func(s *grpc.Server) {
		if s == nil {
			t.Fatal("grpc.Server should not be nil")
		}
	})
}

func TestSetupBroker(t *testing.T) {
	z := NewZepto()
	bp := &testutils.BrokerProviderMock{}
	z.SetupBroker(bp)
	if z.broker == nil {
		t.Fatal("zepto broker should not be nil")
	}
}
