# ── Stage 1: build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api

# ── Stage 2: run ────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/api ./api

# Create dirs that the server writes to at runtime.
# On Render free tier these are ephemeral (data is lost on redeploy/restart).
RUN mkdir -p uploads storage/cvs

EXPOSE 8080

CMD ["./api"]
