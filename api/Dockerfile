FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /api

EXPOSE 8080

CMD [ "/api" ]