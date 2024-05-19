#docker build -t db-app .
#docker run -p 8082:8082 -d db-app

FROM golang:1.22.3-bullseye

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN go build -o /usr/src/app/db-app main.go

CMD ["/usr/src/app/db-app"]