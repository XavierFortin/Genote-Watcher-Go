FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app

FROM base AS prod-deps
WORKDIR /app/src/client
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS build
WORKDIR /app/src/client
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm build

# Build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
COPY --from=build /app/src/client/dist /app/src/client/dist
ENV CGO_ENABLED=0
RUN cd src && go mod download
RUN go build -C src -o /bin/genote-watcher -v -ldflags "-X utils.BuildMode=prod"

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/genote-watcher /bin/genote-watcher/app
EXPOSE 4000
ENTRYPOINT [ "/bin/genote-watcher/app" ]

LABEL Name=genotewatcher
