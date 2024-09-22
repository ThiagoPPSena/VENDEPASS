FROM golang:1.23

# Set the working directory to /app
WORKDIR /app
# Copy the server directory to the container
COPY . ./
# Download the dependencies
RUN go mod download
# Build the server
RUN go build -o /server
# expose the port
EXPOSE 8080

CMD [ "/server" ]