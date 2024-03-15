FROM golang:1.21

WORKDIR /app

COPY . .

#RUN go get -d -v ./...
#
#RUN go build -o /go/bin/app ./cmd
RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "./cmd/main.go"]