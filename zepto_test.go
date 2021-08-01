package zepto

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-zepto/zepto/config"
	"google.golang.org/grpc"
)

func TestNewZepto(t *testing.T) {
	z := NewZepto(config.ZEPTO_TEST_CONFIG)
	z.InitApp()
	if z.logger == nil {
		t.Errorf("Logger should not be nil")
	}
}

func TestSetupHTTP(t *testing.T) {
	z := NewZepto(config.ZEPTO_TEST_CONFIG)
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
	if s.Name != "zepto-test" {
		t.Errorf("Expected name to be equal 'zepto'")
	}
	if s.Version != "1.0.0" {
		t.Errorf("Expected version to be equal 'latest'")
	}
	if s.Age != "0s" {
		t.Errorf("Expected age to be equal '0s'")
	}
}

func TestSetupGRPC(t *testing.T) {
	z := NewZepto(config.ZEPTO_TEST_CONFIG)
	z.SetupGRPC("0.0.0.0:9000", func(s *grpc.Server) {
		if s == nil {
			t.Fatal("grpc.Server should not be nil")
		}
	})
}
