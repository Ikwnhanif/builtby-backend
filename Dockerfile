# --- STAGE 1: Build Engine ---
FROM golang:1.26-alpine AS builder

# Install build essentials untuk library C jika diperlukan
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod dan go.sum duluan supaya layer caching jalan
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build dengan optimasi: 
# -ldflags="-s -w" untuk memperkecil ukuran binary
# CGO_ENABLED=0 agar binary bersifat statis (portable)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o builtby-engine .

# --- STAGE 2: Runtime Environment ---
FROM alpine:latest

# Tambahkan timezone data dan ca-certificates untuk keperluan API/HTTPS
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/builtby-engine .
# Copy .env jika kamu menggunakannya untuk config lokal
COPY --from=builder /app/.env . 

# Buat folder uploads agar tidak error saat mounting volume
RUN mkdir -p uploads

EXPOSE 8080

# Jalankan engine
CMD ["./builtby-engine"]