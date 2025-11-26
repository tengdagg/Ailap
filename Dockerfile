# Build Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Build Backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go

# Final Stage
FROM alpine:latest
WORKDIR /app

# Install dependencies (if any needed for runtime, e.g. ca-certificates)
RUN apk --no-cache add ca-certificates tzdata

# Copy backend binary
COPY --from=backend-builder /app/backend/server .

# Copy frontend static files
COPY --from=frontend-builder /app/frontend/dist ./dist

# Create directory for data (if needed, e.g. sqlite)
RUN mkdir -p data

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]
