FROM golang:1.14 as build

WORKDIR /go/src/go-saas
COPY ./backend /go/src/go-saas
COPY ./go-saas.json /go/src/go-saas

RUN go build -o /go/bin/go-saas

FROM gcr.io/distroless/base

COPY --from=build /go/bin/go-saas /
COPY --from=build /go/src/go-saas/go-saas.json /
CMD ["/go-saas"]