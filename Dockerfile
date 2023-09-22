FROM golang:alpine

# Set the working directory in Docker to /app
WORKDIR /app

# Copy everything from the current directory to the working directory in Docker
COPY . .

# This should pick up the go.mod and download necessary dependencies
RUN go mod download

# Build the application
RUN go build -o kadlab ./main/main.go

# Run the application
CMD ["/app/kadlab"]