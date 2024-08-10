FROM golang:1.22-alpine

WORKDIR /go/src/github.com/dskart/waterfall-engine
COPY . .

RUN go generate ./...
RUN go build .

FROM golang:1.22-alpine

WORKDIR /usr/bin

COPY --from=0 /go/src/github.com/dskart/waterfall-engine/waterfall-engine .
RUN ./waterfall-engine --help > /dev/null

ENTRYPOINT ["/usr/bin/sidekick"]

