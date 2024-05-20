#docker build -t db-app .
#docker run -p 8082:8082 -d db-app

FROM golang:1.22.3-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /usr/src/app/db-app main.go

EXPOSE 8082

CMD ["/usr/src/app/db-app"]