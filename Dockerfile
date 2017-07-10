FROM ubuntu:16.04

RUN apt-get update && apt-get -y install netcat
RUN mkdir -p /smsender/config

COPY bin/smsender /smsender/
COPY config/config.default.yml /
COPY webroot/dist /smsender/webroot/dist/

COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 8080