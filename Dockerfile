FROM golang:1.25.4-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata git

ARG MAIN_PKG=./cmd/server

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/server ${MAIN_PKG}

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/server /app/server
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/app/server"]
