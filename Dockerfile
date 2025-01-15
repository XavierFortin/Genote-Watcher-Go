# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /src/genote-watcher
COPY . .
ENV CGO_ENABLED=0
RUN cd src && go mod download
RUN go build -C src -o /bin/genote-watcher -v -ldflags "-X main.buildMode=prod"

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/genote-watcher /bin/genote-watcher/app
ENTRYPOINT [ "/bin/genote-watcher/app" ]

LABEL Name=genotewatcher