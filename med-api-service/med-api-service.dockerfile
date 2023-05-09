FROM alpine:latest

RUN mkdir /app

COPY medApiServiceApp /app

RUN apk add poppler-utils

EXPOSE 80

CMD ["/app/medApiServiceApp"]