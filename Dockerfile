# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build Flags:
# -w -s: Strip debug information (Reduces size)
# CGO_ENABLED=0: Static binary (No C dependencies)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o aas-edge ./cmd/server/main.go

# Stage 2: Runtime
# "gcr.io/distroless/static-debian12" contains only CA Certs and Timezone data.
FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/aas-edge /aas-edge

EXPOSE 8080
ENTRYPOINT ["/aas-edge"]
