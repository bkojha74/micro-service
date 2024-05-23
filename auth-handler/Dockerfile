#docker build -t auth-app .
#docker run -p 8081:8081 -d auth-app

FROM golang:1.22.3-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN go build -o /usr/src/app/auth-app main.go

CMD ["/usr/src/app/auth-app"]