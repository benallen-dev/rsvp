FROM golang:1.25.7 AS build

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build a fully static binary for Linux
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o rsvp .

FROM golang:1.25.7 as run

# FROM alpine:latest AS run
# RUN apk add --no-cache sqlite-libs

# Set working directory
WORKDIR /app

# Copy binary and config from build stage
COPY --from=build /app/rsvp .
COPY --from=build /app/config.yaml .
COPY --from=build /app/public ./public
COPY --from=build /app/web/templates ./web/templates/

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./rsvp"]

