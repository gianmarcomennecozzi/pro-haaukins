FROM golang:1.13.0
MAINTAINER "Gian Marco Mennecozzi"
WORKDIR /haaukins
COPY . .
ENV DB_USER ${DB_USER}
ENV DB_PASSWORD ${DB_PASSWORD}
ENV DB_NAME ${DB_NAME}
ENV DB_HOST ${DB_HOST}
ENV DB_PORT ${DB_PORT}
RUN go build -o server .
EXPOSE 50051
CMD ["./server"]
