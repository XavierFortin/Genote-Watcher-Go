FROM node:20-slim AS client-build
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app/src/client

# Only copy package files first to leverage caching
COPY src/client/package.json src/client/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# Then copy source and build
COPY src/client ./
RUN pnpm build

FROM golang:alpine AS go-build
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY . .
COPY --from=client-build /app/src/client/dist /app/src/client/dist
RUN go build -C src -o /bin/genote-watcher -v -ldflags "-X 'genote-watcher/utils.BuildMode=prod'"

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=go-build /bin/genote-watcher /bin/app
EXPOSE 4000
ENTRYPOINT ["/bin/app", "--port", "4000"]
