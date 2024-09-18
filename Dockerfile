# Stage 1: Build aplikasi
FROM golang:alpine AS builder

# Install git dan build-essential untuk mendukung proses build aplikasi
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Salin file go mod dan sum
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Salin kode sumber
COPY . .

# Build aplikasi
RUN go build -o main

# Stage 2: Image minimal untuk menjalankan aplikasi
FROM alpine:latest

# Install lib yang mungkin diperlukan oleh aplikasi
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Salin binary dari tahap build
COPY --from=builder /app/main .

# Setting Permission
RUN chmod +x main

# Jalankan aplikasi
CMD ["./main", "-c", "1"]