# Start with the official Golang image
FROM golang:1.19 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o web-app ./cmd/web

# Use a small base image to create a final container
FROM alpine:3.14

# Set the current working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/web-app /app/web-app

# Expose port 8080 (or whatever port your app runs on)
EXPOSE 8080

# Command to run the app
CMD ["/app/web-app"]
