FROM golang:1.23

# Set the working directory to /app
WORKDIR /app
# Copy the client directory to the container
COPY . ./
# Download the dependencies
RUN go mod download
# Build the server
RUN go build -o /client

# Set environment variables
ENV SERVER_ADDRESS=server
ENV SERVER_PORT=8080

CMD [ "/client" ]