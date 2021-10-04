FROM golang:1.17.1-buster as builder

# Create and change to the src directory.
WORKDIR /src

# Copy local code to the container image.
COPY . ./

# Get dependencies
RUN go get -d -v ./...

# Build the binary.
RUN go build -v -o app

# Use the distroless image for a lean and secure container.
FROM gcr.io/distroless/base

# Copy the binary to the production image from the builder stage.
COPY --from=builder /src/app /

# Run the web service on container startup.
CMD ["/app"]