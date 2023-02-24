#
# Cloudmanic Backup Dockerfile
#

FROM golang:1.18-alpine as builder

# Set the working directory
WORKDIR /app

# Copy the source code to the container
COPY . .

# Build the binary executable
RUN go mod tidy
RUN go build -o cloudmanic-backup .

# Final image
FROM alpine:3.7

# Install Helpful Programs
RUN apk --no-cache add bash mysql-client

# Create a user to run the app
RUN adduser -D -u 1000 backup

# Copy the binary executable
COPY --from=builder /app/cloudmanic-backup /usr/local/bin/

# Set the user
USER backup

# Set the working directory
WORKDIR /usr/local/bin

# Start the app
CMD ["./cloudmanic-backup"]