package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadijannat/go-aas-edge/internal/ai"
	"github.com/hadijannat/go-aas-edge/internal/model"
	"github.com/hadijannat/go-aas-edge/internal/physics"
	"github.com/joho/godotenv"
)

func main() {
	// Load env for local dev (ignore error in Docker)
	_ = godotenv.Load()

	// 1. Initialize Domain
	twin := model.NewTwinManager()
	aiAgent := ai.NewAgent()

	// 2. Start Background Simulation
	physics.Run(twin)
	log.Println("⚙️  Physics Simulation Started...")

	if aiAgent.IsConfigured() {
		log.Println("🤖 AI Agent: OpenAI configured")
	} else {
		log.Println("🤖 AI Agent: Running in mock mode (set OPENAI_API_KEY for live AI)")
	}

	r := setupRouter(twin, aiAgent)

	log.Println("🚀 Go-AAS-Edge running on :8080")
	log.Println("   - GET  /health    → Health check")
	log.Println("   - GET  /aas       → AAS V3 JSON")
	log.Println("   - GET  /telemetry → Live sensor data")
	log.Println("   - POST /ask       → Chat with the asset")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// setupRouter wires all HTTP endpoints. It is separated from main to enable
// integration testing without starting the long-running server loop.
func setupRouter(twin *model.TwinManager, aiAgent *ai.Agent) *gin.Engine {
	if gin.Mode() == gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "go-aas-edge",
		})
	})

	// Endpoint A: Standard AAS V3 JSON
	r.GET("/aas", func(c *gin.Context) {
		data, err := twin.SerializeToJSON()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	})

	// Endpoint B: Live Telemetry Snapshot
	r.GET("/telemetry", func(c *gin.Context) {
		snapshot := twin.Snapshot()
		c.JSON(http.StatusOK, snapshot)
	})

	// Endpoint C: Cognitive Chat
	r.POST("/ask", func(c *gin.Context) {
		var req struct {
			Question string `json:"question" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "question field is required"})
			return
		}

		// RAG: Get Snapshot -> Call AI
		snapshot := twin.Snapshot()
		answer, err := aiAgent.Diagnose(c.Request.Context(), snapshot, req.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"telemetry": snapshot,
			"reply":     answer,
		})
	})

	return r
}
