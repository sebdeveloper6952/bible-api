# Stage 1: Build Go binary
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags='-s -w' -o /api ./cmd/api

# Stage 2: Runtime (distroless includes CA certs for outbound TLS)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
CMD ["serve"]
