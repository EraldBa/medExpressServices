FROM alpine:latest

WORKDIR /app

COPY medApiServiceApp .

RUN apk add poppler-utils

EXPOSE 80

ENTRYPOINT ["./medApiServiceApp"]