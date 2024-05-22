#docker build -t db-app .
#docker run -p 8083:8083 -d file-app

FROM golang:1.22.3-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /usr/src/app/file-app main.go

EXPOSE 8083

CMD ["/usr/src/app/file-app"]