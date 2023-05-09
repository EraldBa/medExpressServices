FROM alpine:latest

RUN mkdir /app

COPY medScraperServiceApp /app

EXPOSE 80

CMD ["/app/medScraperServiceApp"]