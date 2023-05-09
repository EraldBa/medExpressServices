FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

EXPOSE 80

CMD ["/app/brokerApp"]