FROM node:20-alpine AS ui-build

WORKDIR /go/src/github.com/dskart/waterfall-engine/ui
COPY ./ui .

RUN mkdir -p ./public/static
RUN wget https://unpkg.com/htmx.org/dist/htmx.min.js -O ./public/static/htmx.min.js
RUN npm install
RUN npm run build


FROM golang:1.22-alpine AS build

WORKDIR /go/src/github.com/dskart/waterfall-engine
COPY . .
COPY --from=ui-build /go/src/github.com/dskart/waterfall-engine/ui ./ui

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go generate ./...
RUN templ generate
RUN go build .


FROM golang:1.22-alpine

WORKDIR /usr/bin

COPY --from=build /go/src/github.com/dskart/waterfall-engine/waterfall-engine .
RUN ./waterfall-engine --help > /dev/null

ENTRYPOINT ["/usr/bin/waterfall-engine"]

