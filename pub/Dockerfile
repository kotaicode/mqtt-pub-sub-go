FROM golang:1.13
ENV GO111MODULE=on

ENV APP_ROOT /app
WORKDIR $APP_ROOT

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o publisher .

CMD ["./publisher"]
EXPOSE 8080
