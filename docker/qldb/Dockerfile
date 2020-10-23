# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Andrii Tarykin <ksbomj@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY services/qldb/go.mod services/qldb/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY services/qldb/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main app/qldb-api/main.go

######## Start a new stage from scratch #######
FROM alpine:latest

ENV QLDB_DATABASE="obada"
ENV QLDB_REGION="eu-central-1"
ENV QLDB_KEY=""
ENV QLDB_SECRET=""

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
