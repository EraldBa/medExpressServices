FROM alpine:latest

WORKDIR /app

COPY medScraperServiceApp .

EXPOSE 80

ENTRYPOINT ["./medScraperServiceApp"]