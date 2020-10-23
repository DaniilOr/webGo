FROM alpine:3.7

ADD bank /app/

ENTRYPOINT ["/app/bank"]

EXPOSE 9999