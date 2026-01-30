# Stage 1: Build
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app \
    -ldflags "-X 'myapp/internal/config.Version=${VERSION}' -s -w" \
    ./cmd/app

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /app/app /app/app
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/app/app"]