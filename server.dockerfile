FROM alpine:3.10.2

WORKDIR /app
ADD myapp /app/

ENTRYPOINT ["./myapp"]