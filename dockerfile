FROM nats:latest

COPY nats-server.conf /etc/nats/nats-server.conf

EXPOSE 4222

EXPOSE 8222

CMD ["-c", "/etc/nats/nats-server.conf" ]

