FROM golang:1.21.3

WORKDIR /src

# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /server

EXPOSE 80
CMD ["/server"]
