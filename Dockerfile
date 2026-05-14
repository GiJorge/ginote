# --- STAGE 1: Build Frontend ---
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# --- STAGE 2: Build Backend ---
FROM golang:1.26-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# We build the binary here
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# --- STAGE 3: Final Combined Image ---
FROM alpine:3.19
RUN apk add --no-cache ca-certificates libc6-compat
WORKDIR /root/

# Copy the Go binary from Stage 2
COPY --from=backend-builder /app/main .

# Copy the Vue static files from Stage 1 into a 'dist' folder
COPY --from=frontend-builder /app/dist ./dist

# Create folder for SQLite database
RUN mkdir -p /root/data

EXPOSE 8389
CMD ["./main"]