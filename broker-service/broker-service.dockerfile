FROM alpine:latest

WORKDIR /app

COPY brokerApp .

EXPOSE 80

ENTRYPOINT ["./brokerApp"]