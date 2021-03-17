# 2-stage build starts from golang base image
FROM golang:alpine as builder

LABEL maintainer="BretH bharris183@gmail.com"

## git is needed for fetching dependencies
RUN apk update && apk add --no-cache git

# Working directory in container
WORKDIR /app

# Copy go mod and go sum
COPY go.mod .
COPY go.sum .

# Download all dependencies
RUN go mod download

# Copy source to container
COPY . .

# Build the go app
RUN CG0_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

# Second stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary and .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8010
EXPOSE 8010

# Run the executable
CMD ["./main"]

