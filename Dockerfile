FROM golang:1.21-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy only go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"] 