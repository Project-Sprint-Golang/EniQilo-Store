# Use the official Golang image as a base
FROM golang:1.22

# Set the current working directory inside the container
WORKDIR /app

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the local package files to the container's workspace
COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o eniqilo .

# # Build the Go app
# RUN go build -o eniqilo .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./eniqilo"]