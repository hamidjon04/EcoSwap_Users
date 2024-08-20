# Stage 1: Build stage
FROM golang:1.22.2 AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the application
COPY . .
# COPY config/model.conf .
COPY .env .

# Optionally copy the .env file if it's needed
#COPY .env .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../myapp .

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app/

RUN apk --no-cache add ca-certificates

RUN mkdir -p /app/logs
RUN touch /app/logs/app.log

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .
# Copy the configuration files
# COPY --from=builder /app/config/model.conf ./config/
#COPY --from=builder /app/config/policy.csv ./config/


# Optionally copy the .env file if it's needed
COPY --from=builder /app/.env .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]