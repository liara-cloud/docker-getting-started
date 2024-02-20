# Use the official Golang image as the base image
FROM golang:latest
        
# Set the working directory inside the container
WORKDIR /app

# Copy the local code to the container
COPY . .

# Download Go modules
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose port 8080 to the outside world
# EXPOSE 8080

# Command to run the executable
CMD ["./main"]