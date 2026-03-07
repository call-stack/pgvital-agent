FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /pgvitals-agent ./cmd/pgvitals-agent

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /pgvitals-agent /usr/local/bin/pgvitals-agent
ENTRYPOINT ["pgvitals-agent"]
