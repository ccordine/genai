# Use the official Go image as the base image
FROM golang:1.23.6-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o genai .

# Set the entry point to the genai script
ENTRYPOINT ["./genai"]

