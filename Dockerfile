FROM golang:1.19-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

RUN go run github.com/prisma/prisma-client-go generate

# Build the binary.
RUN go build -v -o server



# Run the web service on container startup.
CMD ["/app/server"]
