# -- mirror args (overridden by install.sh when user picks China mirrors)
ARG NPM_REGISTRY=https://registry.npmjs.org
ARG GOPROXY=https://proxy.golang.org,direct
ARG ALPINE_MIRROR=https://dl-cdn.alpinelinux.org/alpine

# Stage 1: Build frontend
FROM node:20-alpine AS frontend
ARG NPM_REGISTRY
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm config set registry ${NPM_REGISTRY} && npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.26-alpine AS backend
ARG GOPROXY
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go env -w GOPROXY=${GOPROXY} && go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o tunnel-manager .

# Stage 3: Final image
FROM alpine:3.19
ARG ALPINE_MIRROR
RUN sed -i "s|https://dl-cdn.alpinelinux.org/alpine|${ALPINE_MIRROR%/}|g" /etc/apk/repositories && \
    apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/tunnel-manager .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN mkdir -p data

EXPOSE 8080
ENV PORT=8080
ENV STATIC_DIR=frontend/dist
ENV STORE_PATH=data/config.json

CMD ["./tunnel-manager"]
