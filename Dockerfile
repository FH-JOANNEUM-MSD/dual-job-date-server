# ==========================================
# Stage 1: Builder (Kompiliert dein Go-Backend)
# ==========================================
# Nimmt die aktuellste Go-Version auf Alpine-Basis
FROM golang:alpine AS builder

WORKDIR /app

# Dependencies cachen (extrem wichtig für schnelle Builds!)
COPY go.mod go.sum ./
RUN go mod download

# Restlichen Code kopieren
COPY . .

# HIER IST DIE ÄNDERUNG: Wir zeigen Go genau, wo die main() liegt!
RUN CGO_ENABLED=0 GOOS=linux go build -o server_binary ./cmd/server


# ==========================================
# Stage 2: Finales Image (Schlank und sicher)
# ==========================================
FROM alpine:latest

WORKDIR /app

# Wir holen uns NUR die fertige Binary aus Stage 1
COPY --from=builder /app/server_binary .

# Zeitzonen hinzufügen (wichtig für deine Event-Dates im Code!)
RUN apk --no-cache add tzdata

# Port dokumentieren (falls dein Server z.B. auf 8080 lauscht)
EXPOSE 8080

# Server starten
CMD ["./server_binary"]