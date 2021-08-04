# Dockerfile References: https://docs.docker.com/engine/reference/builder/
FROM alpine:edge AS builder

RUN apk update && apk upgrade

RUN apk add --no-cache go gcc g++ musl-dev

WORKDIR /app

COPY src/ .

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags '-extldflags "-static"' --tags "json1 fts5 secure_delete" -a -o main app/node/main.go

######## Start a new stage from scratch #######
FROM alpine:latest

ENV QLDB_DATABASE="obada"
ENV AWS_REGION="us-east-1"
ENV AWS_KEY=""
ENV AWS_SECRET=""

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]