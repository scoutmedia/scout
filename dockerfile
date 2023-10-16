# Pull latest image of golang
FROM golang:1.21-alpine
# Add bash
RUN apk update && apk add bash
# Create new directory inside of container
WORKDIR /scout
# Copy all files to destination folder
COPY . /scout
# Download go mod files
RUN go mod download
# Execute go run main.go 
RUN go build -o main
# Run main.go file inside workingdir
CMD ["/scout/main"]