FROM golang:1.21.3

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY genproto genproto
COPY *.go ./

RUN go build -o /server

EXPOSE 80
CMD ["/server"]
