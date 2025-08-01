# Use the official Golang image
FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go module files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install air for hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy the rest of the application source code
COPY . .

# Copy the entrypoint script and make it executable
COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh

# Set the entrypoint to our script
ENTRYPOINT ["./entrypoint.sh"]

# The command to run with the entrypoint (air)
CMD ["air", "-c", ".air.toml"]