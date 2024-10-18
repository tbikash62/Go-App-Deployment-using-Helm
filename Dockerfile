# Stage 1: Build the Go binary  
FROM golang:1.23.2-alpine3.20 AS buildstage  

# Set the working directory  
WORKDIR /app  
  
# Copy the Go module files and download dependencies  
RUN go mod init app
RUN go mod download
RUN go mod tidy  
  
# Copy the rest of the application code  
COPY main.go .  
  
# Build the Go binary  
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .  
  
# Stage 2: Create a minimal image  
FROM alpine:latest  

# Added as started getting error in helm!
RUN apk add --no-cache bash 
  
# Copy the built Go binary from the build stage  
COPY --from=buildstage /app/app /app/app

# Create a non-root user  
RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app  
RUN chmod +x /app 
  
# Set the user to the non-root user  
USER appuser 
  
# Set the PORT environment variable with a default value  
ENV PORT=8080
EXPOSE ${PORT}

# Command to run the application
ENTRYPOINT ["/bin/sh", "-c", "./app/app ${PORT}"] 