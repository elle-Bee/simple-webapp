FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
#RUN go build -v -o /usr/local/bin/app ./...

RUN go build -o main .

EXPOSE 8081

#RUN chmod +x main

CMD ["./main"]
