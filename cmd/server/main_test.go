package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hadijannat/go-aas-edge/internal/ai"
	"github.com/hadijannat/go-aas-edge/internal/model"
)

// Build a fresh router for every test to avoid state leakage between cases.
func newTestRouter() (*gin.Engine, *model.TwinManager) {
	gin.SetMode(gin.TestMode)
	os.Unsetenv("OPENAI_API_KEY")
	twin := model.NewTwinManager()
	aiAgent := ai.NewAgent()

	return setupRouter(twin, aiAgent), twin
}

func TestHealthEndpoint(t *testing.T) {
	router, _ := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if payload["status"] != "healthy" || payload["service"] != "go-aas-edge" {
		t.Fatalf("unexpected health payload: %#v", payload)
	}
}

func TestTelemetryEndpoint(t *testing.T) {
	router, twin := newTestRouter()
	twin.UpdateState("42.00", "1234")

	req := httptest.NewRequest(http.MethodGet, "/telemetry", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode telemetry response: %v", err)
	}

	if payload["Temperature"] != "42.00" || payload["RPM"] != "1234" {
		t.Fatalf("unexpected telemetry payload: %#v", payload)
	}
}

func TestAskEndpointUsesTelemetry(t *testing.T) {
	router, twin := newTestRouter()
	twin.UpdateState("33.30", "987")

	body := bytes.NewBufferString(`{"question": "How are you?"}`)
	req := httptest.NewRequest(http.MethodPost, "/ask", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload struct {
		Telemetry map[string]string `json:"telemetry"`
		Reply     string            `json:"reply"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode ask response: %v", err)
	}

	if payload.Telemetry["Temperature"] != "33.30" || payload.Telemetry["RPM"] != "987" {
		t.Fatalf("telemetry mismatch in ask response: %#v", payload.Telemetry)
	}

	if payload.Reply == "" {
		t.Fatalf("expected a reply message, got empty string")
	}
}

func TestAASSerializationEndpoint(t *testing.T) {
	router, _ := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/aas", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode aas response: %v", err)
	}

	shells, ok := payload["assetAdministrationShells"].([]any)
	if !ok || len(shells) == 0 {
		t.Fatalf("expected assetAdministrationShells in response")
	}

	firstShell, ok := shells[0].(map[string]any)
	if !ok || firstShell["idShort"] != "SmartMotorTwin" {
		t.Fatalf("unexpected shell content: %#v", firstShell)
	}
}
