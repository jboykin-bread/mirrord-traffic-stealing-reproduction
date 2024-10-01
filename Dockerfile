# Build stage
FROM golang:1.23 AS builder
WORKDIR /workspace

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o color-server

# Final image
FROM gcr.io/distroless/base-debian11 AS final
WORKDIR /app
COPY --from=builder /workspace/color-server .
USER nonroot:nonroot
ENTRYPOINT ["/app/color-server"]