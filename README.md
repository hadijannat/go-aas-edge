# 🏭 Go-AAS-Edge

Lightweight Asset Administration Shell (AAS) server for edge devices with optional AI diagnostics.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.23+">
  <img src="https://img.shields.io/badge/AAS-V3.0-FF6B35?style=for-the-badge" alt="AAS V3.0">
  <img src="https://img.shields.io/badge/Docker-<20MB-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker <20MB">
</p>

## 🎯 Why Go-AAS-Edge?

| Java AAS Server | Go-AAS-Edge |
|-----------------|-------------|
| ~300MB Docker Image | <20MB Docker Image |
| JVM Required | Single Binary |
| High Memory Footprint | <50MB RAM |
| Complex Deployment | `./aas-edge` |

## ✨ Features

- **📦 Small image**: Distroless container without extra tooling
- **⚡ Concurrency**: Goroutines handle requests and simulation in parallel
- **🔧 AAS V3.0 compliant**: Uses the `aas-core-works` SDK
- **🤖 AI diagnostics**: Optional chat endpoint backed by OpenAI when configured
- **🔄 Real-time simulation**: Physics loop updates telemetry values

## 🚀 Quick Start

### Run with Docker

```bash
# Build the image
docker build -t go-aas-edge .

# Run (with AI enabled)
docker run -p 8080:8080 -e OPENAI_API_KEY="sk-..." go-aas-edge

# Run (mock AI mode)
docker run -p 8080:8080 go-aas-edge
```

### Run Locally

```bash
# Install dependencies
go mod tidy

# Run (optional: set OPENAI_API_KEY for live AI)
go run ./cmd/server/main.go
```

## 📡 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/aas` | GET | Full AAS V3.0 JSON structure |
| `/telemetry` | GET | Live sensor values (Temperature, RPM) |
| `/ask` | POST | Chat with your asset |

## 💬 Chat with Your Asset

```bash
curl -X POST http://localhost:8080/ask \
     -H "Content-Type: application/json" \
     -d '{"question": "Are you overheating?"}'
```

**Response:**
```json
{
  "telemetry": {
    "Temperature": "45.30",
    "RPM": "1250"
  },
  "reply": "My current temperature is 45.3°C. This is well within my safe operating range of 85°C. No overheating detected."
}
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Go-AAS-Edge                            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────────┐   │
│  │   Physics   │──▶│  TwinManager│──▶│   REST API      │   │
│  │   Engine    │   │  (AAS V3.0) │   │   (Gin)         │   │
│  └─────────────┘   └─────────────┘   └────────┬────────┘   │
│        │                  │                    │            │
│        │                  ▼                    │            │
│        │           ┌─────────────┐             │            │
│        └──────────▶│  AI Agent   │◀────────────┘            │
│                    │  (OpenAI)   │                          │
│                    └─────────────┘                          │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
go-aas-edge/
├── cmd/server/main.go      # Application entry point
├── internal/
│   ├── model/twin.go       # Thread-safe AAS wrapper
│   ├── physics/engine.go   # Sensor simulation
│   └── ai/agent.go         # AI diagnostics client
├── Dockerfile              # Multi-stage distroless build
├── .env.example            # Configuration template
└── go.mod
```

## 🔒 Security

- **Distroless Base**: No shell, no package manager
- **Static Binary**: No runtime dependencies
- **Thread-Safe**: Mutex-protected state updates

## 📄 License

MIT License - See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [aas-core-works](https://github.com/aas-core-works/aas-core3.0-golang) - AAS V3.0 SDK
- [Gin](https://github.com/gin-gonic/gin) - HTTP framework
- [OpenAI](https://platform.openai.com/) - AI capabilities
