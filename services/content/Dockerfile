FROM golang:1.21.3

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/a-h/templ/cmd/templ@latest
COPY *.templ ./
RUN templ generate

COPY genproto genproto
COPY *.go ./

RUN go build -o /server

EXPOSE 80
CMD ["/server"]
