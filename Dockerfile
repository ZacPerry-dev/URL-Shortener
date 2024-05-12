FROM golang:latest

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@latest

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose the port your application runs on
EXPOSE 8000

# Command to run the application (using air)
CMD ["air", "-c", ".air.toml"]
