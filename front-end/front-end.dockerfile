FROM alpine:latest

RUN mkdir /app

COPY bin/frontendApp /app

CMD [ "/app/frontendApp" ]
