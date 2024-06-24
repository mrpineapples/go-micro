FROM alpine:latest

RUN mkdir /app

COPY bin/mailApp /app
COPY templates /templates

CMD [ "/app/mailApp" ]