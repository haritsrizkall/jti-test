# Stage 1: Build Stage
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary Go application files to the container
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o app main.go

# Stage 2: Production Stage
FROM alpine:latest

ENV TZ=Asia/Bangkok

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the build stage
COPY --from=builder /app/app .
COPY --from=builder /app/.env ./.env
COPY --from=builder /app/views ./views
COPY --from=builder /app/sounds ./sounds

# Expose the port the app runs on
EXPOSE 1323

# Command to run the executable
CMD ["./app"]
