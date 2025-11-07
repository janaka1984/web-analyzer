# ────────────────────────────────
# 1. Test stage
# ────────────────────────────────
FROM golang:1.25-alpine AS tester

WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Run tests with coverage — build fails if tests fail
RUN go test ./... -cover


# ────────────────────────────────
# 2. Build stage
# ────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of your source code (includes templates)
COPY . .

# verify templates exist
RUN ls -R internal/view/templates

# build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api


# ────────────────────────────────
# 3. Runtime stage
# ────────────────────────────────
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# copy built binary
COPY --from=builder /app/api .

# ✅ copy templates for Gin to render HTML
COPY --from=builder /app/internal/view/templates ./internal/view/templates

EXPOSE 8080
ENTRYPOINT ["/app/api"]
