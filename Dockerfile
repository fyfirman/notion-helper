FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN make build

ENTRYPOINT ["/app/notion-helper"]