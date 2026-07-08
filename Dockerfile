# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o tunnel-manager .

# Stage 3: Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend /app/tunnel-manager .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN mkdir -p data

EXPOSE 8080
ENV PORT=8080
ENV STATIC_DIR=frontend/dist
ENV STORE_PATH=data/config.json

CMD ["./tunnel-manager"]
