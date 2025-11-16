FROM golang:1.24-alpine AS modules

RUN apk add --no-cache git

COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:1.24-alpine

RUN apk add --no-cache git

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./integration-test/..."]
