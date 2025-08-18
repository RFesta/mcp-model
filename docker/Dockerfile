# --- builder
FROM golang:1.22 AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /out/app ./cmd/modelo-mcp

# --- runner
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/app /app/app
COPY configs /app/configs
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/app"]
